package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "net/http/pprof" // include pprop
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	"github.com/mweagle/Sparta/interceptor"
	gocf "github.com/mweagle/go-cloudformation"
	"github.com/sirupsen/logrus"
)

const (
	magicErrorValue = "invalidInput"
)

/*
Supported signatures

• func ()
• func () error
• func (TIn), error
• func () (TOut, error)
• func (context.Context) error
• func (context.Context, TIn) error
• func (context.Context) (TOut, error)
• func (context.Context, TIn) (TOut, error)
*/

// Standard AWS λ function

func hellowXRayWorld(ctx context.Context, msg json.RawMessage) (string, error) {
	logger, _ := ctx.Value(sparta.ContextKeyLogger).(*logrus.Logger)
	logger.Info("Regular log message goes here 🎁")

	logger.Trace("Checking magic value. Only logged on lambda error.")
	if strings.Contains(string(msg), magicErrorValue) {
		logger.WithField("rejectedValue", magicErrorValue).
			Debug("This log message will only be logged if the lambda returns an error. You can log a redacted event here.")
		return "", errors.New("This magic value isn't supported")
	}
	return "XRayInterceptor executed normally", nil
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	lambdaFn, _ := sparta.NewAWSLambda("HelloXRayWorld",
		hellowXRayWorld,
		sparta.IAMRoleDefinition{})
	lambdaFn.Options.TracingConfig = &gocf.LambdaFunctionTracingConfig{
		Mode: gocf.String("Active"),
	}

	lambdaFn.Interceptors = interceptor.RegisterXRayInterceptor(lambdaFn.Interceptors,
		interceptor.XRayAll)

	sess := session.Must(session.NewSession())
	awsName, awsNameErr := spartaCF.UserAccountScopedStackName("SpartaXRayInterceptor",
		sess)
	if awsNameErr != nil {
		fmt.Println("Failed to create stack name")
		os.Exit(1)
	}
	// Create the stack...
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFunctions = append(lambdaFunctions, lambdaFn)
	err := sparta.Main(awsName,
		"Simple Sparta application that demonstrates how to support XRay interceptor support",
		lambdaFunctions,
		nil,
		nil)
	if err != nil {
		os.Exit(1)
	}
}
