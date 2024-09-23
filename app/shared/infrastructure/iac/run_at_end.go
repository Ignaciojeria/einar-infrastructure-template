package iac

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func init() {
	ioc.RegistryAtEnd(
		runAtEnd,
		NewPulumiResourceManager)
}
func runAtEnd(rm *PulumiResourceManager) {
	pulumi.Run(func(ctx *pulumi.Context) error {
		return rm.Execute(ctx)
	})
}
