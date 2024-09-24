package bucket

import (
	"iac/app/googlecloud/serviceaccount"
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
		configuration.NewConf,
		serviceaccount.NewOpenObserveSA)
}

type OpenObserveBucket struct {
	Bucket *storage.Bucket
}

func NewOpenObserveBucket(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf,
	sa *serviceaccount.OpenObserveSA) *OpenObserveBucket {
	var openObserveBucket OpenObserveBucket
	rm.Register(
		func(ctx *pulumi.Context) error {
			name := conf.LoadFromSystem("OPENOBSERVE_GCS_BUCKET_NAME")
			bk, err := storage.NewBucket(ctx, name, &storage.BucketArgs{
				Project:                  pulumi.String(conf.GOOGLE_PROJECT_ID),
				Name:                     pulumi.String(name),
				Location:                 pulumi.String("SOUTHAMERICA-WEST1"), // Región de Santiago, Chile
				ForceDestroy:             pulumi.Bool(true),
				UniformBucketLevelAccess: pulumi.Bool(true),
			})
			if err != nil {
				return err
			}
			openObserveBucket.Bucket = bk

			// Add storage admin permissions to the service account
			_, err = storage.NewBucketIAMBinding(ctx, "storageAdminBinding", &storage.BucketIAMBindingArgs{
				Bucket: bk.Name, // specify the bucket name if needed
				Role:   pulumi.String("roles/storage.admin"),
				Members: pulumi.StringArray{
					sa.ServiceAccount.Email.ApplyT(func(email string) (string, error) {
						return "serviceAccount:" + email, nil
					}).(pulumi.StringOutput),
				},
			})
			if err != nil {
				return err
			}

			return nil
		},
	)
	return &openObserveBucket
}
