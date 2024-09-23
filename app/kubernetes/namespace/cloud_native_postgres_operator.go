package namespace

import (
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
		NewCloudNativePostgresOperator,
		iac.NewPulumiResourceManager,
		configuration.NewConf)
}

type CloudNativePostgresOperator struct {
	cloudnativePGOperatorChart *helmv4.Chart
}

func NewCloudNativePostgresOperator(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf) *CloudNativePostgresOperator {
	//kubectl get namespaces
	var cloudNativePostgreSQL CloudNativePostgresOperator
	rm.Register(func(ctx *pulumi.Context) error {
		name := "cloud-native-postgres-operator"
		cnpgSystemNamespace, err := namespace.NewNamespace(ctx,
			name,
			&namespace.NamespaceArgs{
				Metadata: metav1.ObjectMetaArgs{
					ClusterName: pulumi.String(conf.KUBERNETES_CLUSTER_NAME),
					Name:        pulumi.String(name),
				},
			})
		if err != nil {
			return err
		}
		//https://www.pulumi.com/registry/packages/kubernetes/api-docs/helm/v4/chart/
		//https://github.com/cloudnative-pg/charts
		//verify pods : kubectl get pods -n cnpg-system
		cloudnativePGChartID := "cloudnative-pg"
		cloudnativePGOpetarorChart, err := helmv4.NewChart(ctx,
			cloudnativePGChartID,
			&helmv4.ChartArgs{
				Namespace: pulumi.String(name),
				Chart:     pulumi.String(cloudnativePGChartID),
				RepositoryOpts: &helmv4.RepositoryOptsArgs{
					Repo: pulumi.String("https://cloudnative-pg.github.io/charts"),
				},
			}, pulumi.DependsOn([]pulumi.Resource{cnpgSystemNamespace}))
		if err != nil {
			return err
		}
		cloudNativePostgreSQL.cloudnativePGOperatorChart = cloudnativePGOpetarorChart
		return nil
	})
	return &cloudNativePostgreSQL

}
