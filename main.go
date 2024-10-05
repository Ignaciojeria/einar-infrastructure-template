package main

import (
	_ "embed"
	_ "iac/app/openobserve_gcp"
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
