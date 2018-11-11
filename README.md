# SpartaXRayInterceptor
Sparta-based application that shows how to use [Interceptors](https://godoc.org/github.com/mweagle/Sparta#LambdaAWSInfo)
to attach a buffered logger and XRay segment publsiher


1. [Install Go](https://golang.org/doc/install)
1. `go get github.com/mweagle/SpartaXRayInterceptor`
1. `cd ./SpartaXRayInterceptor`
1. `go get -u -v ./...`
1. `go get -u -v github.com/magefile/mage`
1. `mage provision`
1. Visit the AWS Console and test your function!