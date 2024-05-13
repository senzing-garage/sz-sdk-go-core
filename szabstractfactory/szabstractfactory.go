package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/sz"
)

// Szabstractfactory is an implementation of the sz.SzAbstractFactory interface.
type Szabstractfactory struct {
	ConfigId       int64
	InstanceName   string
	Settings       string
	VerboseLogging int64
}

// ----------------------------------------------------------------------------
// sz.SzAbstractFactory interface methods
// ----------------------------------------------------------------------------

/*
TODO: Write description for CreateSzConfig
The CreateSzConfig method...

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzConfig object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (sz.SzConfig, error) {
	result := &szconfig.Szconfig{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzConfigManager
The CreateSzConfigManager method...

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.CreateConfigManager object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (sz.SzConfigManager, error) {
	result := &szconfigmanager.Szconfigmanager{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzDiagnostic
The CreateSzDiagnostic method...

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzDiagnostic object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (sz.SzDiagnostic, error) {
	result := &szdiagnostic.Szdiagnostic{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigId, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzEngine
The CreateSzEngine method...

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzEngine object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (sz.SzEngine, error) {
	result := &szengine.Szengine{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigId, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzProduct
The CreateSzProduct method...

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzProduct object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (sz.SzProduct, error) {
	result := &szproduct.Szproduct{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}
