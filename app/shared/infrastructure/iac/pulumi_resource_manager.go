package iac

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ResourceFunc define el tipo de función que acepta un contexto de Pulumi y retorna un error.
type ResourceFunc func(ctx *pulumi.Context) error

type PulumiResourceManager struct {
	resourceFuncs []ResourceFunc
}

func init() {
	ioc.Registry(NewPulumiResourceManager)
}

// NewResourceManager es el constructor que inicializa un nuevo ResourceManager.
func NewPulumiResourceManager() *PulumiResourceManager {
	return &PulumiResourceManager{
		resourceFuncs: []ResourceFunc{},
	}
}

// Register añade una nueva función de recurso al manager.
func (rm *PulumiResourceManager) Register(f ResourceFunc) {
	rm.resourceFuncs = append(rm.resourceFuncs, f)
}

// Execute ejecuta todas las funciones de recursos en el orden registrado.
func (rm *PulumiResourceManager) execute(ctx *pulumi.Context) error {
	for _, f := range rm.resourceFuncs {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}
