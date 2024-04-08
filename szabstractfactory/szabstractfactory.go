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

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Szconfig is the default implementation of the Szconfig interface.
type Szabstractfactory struct {
	ConfigId     int64
	InstanceName string
	// isTrace        bool
	// logger         logging.LoggingInterface
	// observerOrigin string
	// observers      subject.Subject
	Settings       string
	VerboseLogging int64
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The CreateConfig method... TODO:

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzConfig object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateConfig(ctx context.Context) (sz.SzConfig, error) {
	result := &szconfig.Szconfig{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
The CreateConfigManager method... TODO:

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.CreateConfigManager object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateConfigManager(ctx context.Context) (sz.SzConfigManager, error) {
	result := &szconfigmanager.Szconfigmanager{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}

/*
The CreateDiagnostic method... TODO:

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzDiagnostic object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateDiagnostic(ctx context.Context) (sz.SzDiagnostic, error) {
	result := &szdiagnostic.Szdiagnostic{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging, factory.ConfigId)
	return result, err
}

/*
The CreateEngine method... TODO:

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzEngine object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateEngine(ctx context.Context) (sz.SzEngine, error) {
	result := &szengine.Szengine{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging, factory.ConfigId)
	return result, err
}

/*
The CreateProduct method... TODO:

Input
  - ctx: A context to control lifecycle.

Output
  - An sz.SzProduct object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateProduct(ctx context.Context) (sz.SzProduct, error) {
	result := &szproduct.Szproduct{}
	err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	return result, err
}
