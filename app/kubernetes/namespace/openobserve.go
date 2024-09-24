package namespace

import (
	"iac/app/googlecloud/serviceaccount"
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
		NewOpenObserve,
		iac.NewPulumiResourceManager,
		configuration.NewConf,
		NewCloudNativePostgresOperator,
		serviceaccount.NewOpenObserveSA)
}
func NewOpenObserve(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf,
	opetaror *CloudNativePostgresOperator,
	sa *serviceaccount.OpenObserveSA,
) {
	name := "openobserve"
	rm.Register(func(ctx *pulumi.Context) error {
		//kubectl get namespaces
		ns, err := namespace.NewNamespace(ctx, name, &namespace.NamespaceArgs{
			Metadata: metav1.ObjectMetaArgs{
				ClusterName: pulumi.String(conf.KUBERNETES_CLUSTER_NAME),
				Name:        pulumi.String(name),
			},
		}, pulumi.DependsOn([]pulumi.Resource{opetaror.cloudnativePGOperatorChart, sa.HmacKey}))
		if err != nil {
			return err
		}
		openobserveChartID := "openobserve"
		_, err = helmv4.NewChart(ctx,
			openobserveChartID,
			&helmv4.ChartArgs{
				Namespace: pulumi.String(name),
				Chart:     pulumi.String(openobserveChartID),
				RepositoryOpts: &helmv4.RepositoryOptsArgs{
					Repo: pulumi.String("https://charts.openobserve.ai"),
				},
			}, pulumi.DependsOn([]pulumi.Resource{ns}))
		if err != nil {
			return err
		}
		return nil
	})

}
