package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go-core/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// Szabstractfactory is an implementation of the senzing.SzAbstractFactory interface.
type Szabstractfactory struct {
	ConfigID       int64
	InstanceName   string
	Settings       string
	VerboseLogging int64
}

// ----------------------------------------------------------------------------
// senzing.SzAbstractFactory interface methods
// ----------------------------------------------------------------------------

/*
TODO: Write description for CreateSzConfig
The CreateSzConfig method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzConfig object.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (senzing.SzConfig, error) {
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
  - An senzing.CreateConfigManager object.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
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
  - An senzing.SzDiagnostic object.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	result := &szdiagnostic.Szdiagnostic{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzEngine
The CreateSzEngine method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzEngine object.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (senzing.SzEngine, error) {
	result := &szengine.Szengine{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	return result, err
}

/*
TODO: Write description for CreateSzProduct
The CreateSzProduct method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzProduct object.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	result := &szproduct.Szproduct{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}
