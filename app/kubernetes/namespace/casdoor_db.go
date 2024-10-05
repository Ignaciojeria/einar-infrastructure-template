package namespace

import (
	"encoding/base64"
	"fmt"
	"iac/app/shared/configuration"
	"iac/app/shared/infrastructure/iac"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	namespace "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func init() {
	ioc.Registry(
		NewPostgresCluster,
		iac.NewPulumiResourceManager,
		configuration.NewConf,
		NewCloudNativePostgresOperator)
}

type PostgresCluster struct {
	postgresClusterChart *helmv4.Chart
	ConnectionString     string
}

func NewPostgresCluster(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf,
	operator *CloudNativePostgresOperator) *PostgresCluster {
	// Definir el namespace para el clúster PostgreSQL
	namespaceName := "database"

	var postgresCluster PostgresCluster
	rm.Register(func(ctx *pulumi.Context) error {
		// Evitar que el código se ejecute durante el preview
		if ctx.DryRun() {
			fmt.Println("Previsualización: no se ejecutarán acciones que dependan de recursos creados.")
			return nil
		}

		// Crear el namespace para el clúster PostgreSQL, depende del operador
		_, err := namespace.NewNamespace(ctx,
			namespaceName,
			&namespace.NamespaceArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name: pulumi.String(namespaceName),
				},
			}, pulumi.DependsOn([]pulumi.Resource{operator.cloudnativePGOperatorChart}))
		if err != nil {
			return err
		}

		// Instalar el clúster de PostgreSQL usando el chart cnpg/cluster
		postgresChartID := "cluster"
		postgresClusterChart, err := helmv4.NewChart(ctx,
			postgresChartID,
			&helmv4.ChartArgs{
				Namespace: pulumi.String(namespaceName),
				Chart:     pulumi.String(postgresChartID),
				RepositoryOpts: &helmv4.RepositoryOptsArgs{
					Repo: pulumi.String("https://cloudnative-pg.github.io/charts"),
				},
			}, pulumi.DependsOn([]pulumi.Resource{operator.cloudnativePGOperatorChart}))
		if err != nil {
			return err
		}
		postgresCluster.postgresClusterChart = postgresClusterChart

		// Obtener el Secret que contiene las credenciales de PostgreSQL (solo durante ejecución real)
		secret, err := namespace.GetSecret(ctx,
			"cluster-superuser",
			pulumi.ID("database/cluster-superuser"), nil,
			pulumi.DependsOn([]pulumi.Resource{postgresClusterChart}))
		if err != nil {
			return err
		}

		// Extraer los datos del Secret
		secret.Data.ApplyT(func(data map[string]string) (string, error) {
			// Decodificar las credenciales
			dbUser, err := decodeBase64(data["username"])
			if err != nil {
				return "", err
			}

			dbPassword, err := decodeBase64(data["password"])
			if err != nil {
				return "", err
			}

			dbHost := "cluster-rw.database.svc.cluster.local" // Cambiar host a cluster-rw en el namespace database
			dbPort := "5432"                                  // Asignar puerto fijo

			// Generar el string de conexión
			connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/database", dbUser, dbPassword, dbHost, dbPort)

			// Imprimir el string de conexión en el formato adecuado
			fmt.Println("Connection String:", connectionString)

			// Almacenar el connection string en el struct para uso posterior
			postgresCluster.ConnectionString = connectionString
			return connectionString, nil
		})

		return nil
	})
	return &postgresCluster
}

// Helper function to decode base64
func decodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
