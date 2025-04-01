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
	ConfigID                     int64
	InstanceName                 string
	Settings                     string
	VerboseLogging               int64
	isSzconfigInitialized        bool
	isSzconfigmanagerInitialized bool
	isSzdiagnosticInitialized    bool
	isSzengineInitialized        bool
	isSzproductInitialized       bool
}

// ----------------------------------------------------------------------------
// senzing.SzAbstractFactory interface methods
// ----------------------------------------------------------------------------

/*
Method CreateConfigManager returns an SzConfigManager object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfigManager object.
*/
func (factory *Szabstractfactory) CreateConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	var err error

	result := &szconfigmanager.Szconfigmanager{}

	if !factory.isSzconfigmanagerInitialized {
		err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
		if err == nil {
			factory.isSzconfigmanagerInitialized = true
		}
	}

	return result, err
}

/*
Method CreateDiagnostic returns an SzDiagnostic object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzDiagnostic object.
*/
func (factory *Szabstractfactory) CreateDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	var err error

	result := &szdiagnostic.Szdiagnostic{}

	if !factory.isSzdiagnosticInitialized {
		err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
		if err == nil {
			factory.isSzdiagnosticInitialized = true
		}
	}

	return result, err
}

/*
Method CreateEngine returns an SzEngine object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzEngine object.
*/
func (factory *Szabstractfactory) CreateEngine(ctx context.Context) (senzing.SzEngine, error) {
	var err error

	result := &szengine.Szengine{}

	if !factory.isSzengineInitialized {
		err := result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
		if err == nil {
			factory.isSzengineInitialized = true
		}
	}

	return result, err
}

/*
Method CreateProduct returns an SzProduct object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzProduct object.
*/
func (factory *Szabstractfactory) CreateProduct(ctx context.Context) (senzing.SzProduct, error) {
	var err error

	result := &szproduct.Szproduct{}

	if !factory.isSzproductInitialized {
		err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
		if err == nil {
			factory.isSzproductInitialized = true
		}
	}

	return result, err
}

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Destroy(ctx context.Context) error {
	var err error

	if factory.isSzconfigInitialized {
		szConfig := &szconfig.Szconfig{}

		err = szConfig.Destroy(ctx)
		if err != nil {
			return err
		}

		factory.isSzconfigInitialized = false
	}

	if factory.isSzconfigmanagerInitialized {
		szConfigmanager := &szconfigmanager.Szconfigmanager{}

		err = szConfigmanager.Destroy(ctx)
		if err != nil {
			return err
		}

		factory.isSzconfigmanagerInitialized = false
	}

	if factory.isSzdiagnosticInitialized {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}

		err = szDiagnostic.Destroy(ctx)
		if err != nil {
			return err
		}

		factory.isSzdiagnosticInitialized = false
	}

	if factory.isSzengineInitialized {
		szEngine := &szengine.Szengine{}

		err = szEngine.Destroy(ctx)
		if err != nil {
			return err
		}

		factory.isSzengineInitialized = false
	}

	if factory.isSzproductInitialized {
		szProduct := &szproduct.Szproduct{}

		err = szProduct.Destroy(ctx)
		if err != nil {
			return err
		}

		factory.isSzproductInitialized = false
	}

	return err
}

/*
Method Reinitialize re-initializes the Senzing objects created by the AbstractFactory with a specific Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (factory *Szabstractfactory) Reinitialize(ctx context.Context, configID int64) error {
	var err error

	factory.ConfigID = configID
	if factory.isSzdiagnosticInitialized {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}

		err = szDiagnostic.Reinitialize(ctx, configID)
		if err != nil {
			return err
		}
	}

	if factory.isSzengineInitialized {
		szEngine := &szengine.Szengine{}

		err = szEngine.Reinitialize(ctx, configID)
		if err != nil {
			return err
		}
	}

	return err
}
