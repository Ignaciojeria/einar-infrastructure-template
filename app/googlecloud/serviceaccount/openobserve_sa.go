package serviceaccount

import (
	"iac/app/shared/configuration"
	"iac/app/shared/infrastructure/iac"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func init() {
	ioc.Registry(
		NewOpenObserveSA,
		iac.NewPulumiResourceManager,
		configuration.NewConf)
}

type OpenObserveSA struct {
	HmacKey *storage.HmacKey
}

func NewOpenObserveSA(
	rm *iac.PulumiResourceManager,
	conf configuration.Conf) *OpenObserveSA {
	var openObserveSA OpenObserveSA
	rm.Register(
		func(ctx *pulumi.Context) error {
			openObserveServiceAccountName := "openobserve-sa"
			serviceAccount, err := serviceaccount.NewAccount(
				ctx,
				openObserveServiceAccountName,
				&serviceaccount.AccountArgs{
					Project:                   pulumi.String(conf.GOOGLE_PROJECT_ID),
					Description:               pulumi.String("open observe GCS service account"),
					AccountId:                 pulumi.String(openObserveServiceAccountName),
					DisplayName:               pulumi.String("Open Observe GCS Account"),
					CreateIgnoreAlreadyExists: pulumi.BoolPtr(true),
				},
			)
			if err != nil {
				return err
			}
			// Create the HMAC key for the associated service account
			key, err := storage.NewHmacKey(ctx, "key", &storage.HmacKeyArgs{
				Project:             pulumi.String(conf.GOOGLE_PROJECT_ID),
				ServiceAccountEmail: serviceAccount.Email,
			})
			if err != nil {
				return err
			}
			openObserveSA.HmacKey = key
			return nil
		},
	)
	return &openObserveSA
}
