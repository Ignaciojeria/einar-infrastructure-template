package configuration

import (
	"iac/app/shared/infrastructure/iac"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Conf struct {
	PORT                        string `required:"true"`
	PROJECT_NAME                string `required:"true"`
	ENVIRONMENT                 string `required:"true"`
	KUBERNETES_CLUSTER_NAME     string `required:"true"`
	GOOGLE_PROJECT_ID           string `required:"true"`
	OPENOBSERVE_GCS_BUCKET_NAME string `required:"true"`
}

func init() {
	ioc.Registry(NewConf, iac.NewPulumiResourceManager)
}

// NewConf accede a las variables definidas en el YAML de Pulumi y las retorna como configuración
func NewConf(rm *iac.PulumiResourceManager) Conf {
	// Obtener las configuraciones desde Pulumi YAML usando ctx.GetConfig()
	var conf Conf
	rm.Register(func(ctx *pulumi.Context) error {
		port, exists := ctx.GetConfig("einar:port")
		if !exists || port == "" {
			port = "8080" // Valor por defecto si no está definida
		}

		projectName, _ := ctx.GetConfig("einar:project_name")
		environment, _ := ctx.GetConfig("einar:environment")
		kubernetesClusterName, _ := ctx.GetConfig("einar:kubernetes_cluster_name")
		googleProjectID, _ := ctx.GetConfig("einar:google_project_id")
		openObserveGcsBucketName, _ := ctx.GetConfig("einar:openobserve_gcs_bucket_name")

		// Crear la configuración
		config := Conf{
			PORT:                        port,
			PROJECT_NAME:                projectName,
			ENVIRONMENT:                 environment,
			KUBERNETES_CLUSTER_NAME:     kubernetesClusterName,
			GOOGLE_PROJECT_ID:           googleProjectID,          // Puede ser vacío
			OPENOBSERVE_GCS_BUCKET_NAME: openObserveGcsBucketName, // Nueva variable
		}

		// Validar la configuración (puedes implementar validaciones adicionales aquí)
		cnf, err := validateConfig(config)
		if err != nil {
			return err
		}
		conf = cnf
		return nil
	})
	return conf
}
