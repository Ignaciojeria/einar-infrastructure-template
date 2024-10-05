package namespace

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func init() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Crear el namespace para Casdoor
		_, err := v1.NewNamespace(ctx, "casdoor", &v1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("casdoor"),
			},
		})
		if err != nil {
			return err
		}

		// Instalar el Casdoor Helm Chart desde el OCI público en Docker Hub
		_, err = helm.NewChart(ctx, "casdoor", helm.ChartArgs{
			Chart:   pulumi.String("casdoor-helm-charts"),
			Version: pulumi.String("v1.714.0"),
			FetchArgs: helm.FetchArgs{
				Repo: pulumi.String("oci://registry-1.docker.io/casbin/casdoor-helm-charts"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
