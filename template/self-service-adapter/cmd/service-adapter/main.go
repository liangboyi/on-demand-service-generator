package main

import (
	"os"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/pivotal-cf/self-service-adapter/adapter"
)

func main() {
	manifestGenerator := new(adapter.ManifestGenerator)
	binder := new(adapter.Binder)
	serviceadapter.HandleCommandLineInvocation(os.Args, manifestGenerator, binder, nil)
}
