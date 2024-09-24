package namespace

import (
	"iac/app/googlecloud/bucket"
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
		serviceaccount.NewOpenObserveSA,
		bucket.NewOpenObserveBucket)
}
func NewOpenObserve(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf,
	opetaror *CloudNativePostgresOperator,
	sa *serviceaccount.OpenObserveSA,
	bk *bucket.OpenObserveBucket,
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
				Values: pulumi.Map{
					"auth": pulumi.Map{
						"ZO_S3_ACCESS_KEY": sa.HmacKey.AccessId,
						"ZO_S3_SECRET_KEY": sa.HmacKey.Secret,
					},
					"config": pulumi.Map{
						"ZO_S3_SERVER_URL":         pulumi.String("https://storage.googleapis.com"),
						"ZO_S3_BUCKET_NAME":        bk.Bucket.Name,
						"ZO_S3_REGION_NAME":        pulumi.String("auto"),
						"ZO_S3_PROVIDER":           pulumi.String("s3"),
						"ZO_S3_FEATURE_HTTP1_ONLY": pulumi.String("true"),
					},
				},
			}, pulumi.DependsOn([]pulumi.Resource{ns}))
		if err != nil {
			return err
		}
		return nil
	})

}
