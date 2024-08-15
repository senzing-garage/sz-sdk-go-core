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

/*
Szabstractfactory is an implementation of the [senzing.SzAbstractFactory] interface.

[senzing.SzAbstractFactory]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go/senzing#SzAbstractFactory
*/
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
The CreateSzConfig method returns an SzConfig object
implemented to use the Senzing native C binary, libG2.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfig object.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	result := &szconfig.Szconfig{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
The CreateSzConfigManager method returns an SzConfigManager object
implemented to use the Senzing native C binary, libG2.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfigManager object.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	result := &szconfigmanager.Szconfigmanager{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
The CreateSzDiagnostic method returns an SzDiagnostic object
implemented to use the Senzing native C binary, libG2.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzDiagnostic object.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	result := &szdiagnostic.Szdiagnostic{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	return result, err
}

/*
The CreateSzEngine method returns an SzEngine object
implemented to use the Senzing native C binary, libG2.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzEngine object.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (senzing.SzEngine, error) {
	result := &szengine.Szengine{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	return result, err
}

/*
The CreateSzProduct method returns an SzProduct object
implemented to use the Senzing native C binary, libG2.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzProduct object.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	result := &szproduct.Szproduct{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}
