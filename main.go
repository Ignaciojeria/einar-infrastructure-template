package main

import (
	_ "embed"
	_ "iac/app/googlecloud/bucket"
	_ "iac/app/googlecloud/serviceaccount"
	_ "iac/app/kubernetes/namespace"
	_ "iac/app/shared/configuration"
	"iac/app/shared/constants"
	"log"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

//go:embed .version
var version string

func main() {
	os.Setenv(constants.Version, version)
	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal(err)
	}
}
