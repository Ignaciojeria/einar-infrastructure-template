package bucket

import (
	"iac/app/shared/configuration"
	"iac/app/shared/infrastructure/iac"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func init() {
	ioc.Registry(
		NewOpenObserveBucket,
		iac.NewPulumiResourceManager,
		configuration.NewConf)
}
func NewOpenObserveBucket(rm *iac.PulumiResourceManager, conf configuration.Conf) {
	rm.Register(
		func(ctx *pulumi.Context) error {
			name := conf.LoadFromSystem("OPENOBSERVE_GCS_BUCKET_NAME")
			_, err := storage.NewBucket(ctx, name, &storage.BucketArgs{
				Project:                  pulumi.String(conf.GOOGLE_PROJECT_ID),
				Name:                     pulumi.String(name),
				Location:                 pulumi.String("SOUTHAMERICA-WEST1"), // Región de Santiago, Chile
				ForceDestroy:             pulumi.Bool(true),
				UniformBucketLevelAccess: pulumi.Bool(true),
			})
			if err != nil {
				return err
			}
			return nil
		},
	)
}
