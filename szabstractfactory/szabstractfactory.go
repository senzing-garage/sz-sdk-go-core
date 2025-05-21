package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/go-helpers/wraperror"
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
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	factory.isSzconfigmanagerInitialized = true

	return result, wraperror.Errorf(err, wraperror.NoMessage)
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
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	factory.isSzdiagnosticInitialized = true

	return result, wraperror.Errorf(err, wraperror.NoMessage)
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
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	factory.isSzengineInitialized = true

	return result, wraperror.Errorf(err, wraperror.NoMessage)
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
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	factory.isSzproductInitialized = true

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Destroy(ctx context.Context) error {
	var err error

	err = factory.destroySzConfigmanager(ctx)
	if err != nil {
		return wraperror.Errorf(err, "SzConfigmanager")
	}

	err = factory.destroySzDiagnostic(ctx)
	if err != nil {
		return wraperror.Errorf(err, "SzDiagnostic")
	}

	err = factory.destroySzEngine(ctx)
	if err != nil {
		return wraperror.Errorf(err, "SzEngine")
	}

	err = factory.destroySzProduct(ctx)
	if err != nil {
		return wraperror.Errorf(err, "SzProduct")
	}

	return nil
}

/*
Method Reinitialize re-initializes the Senzing objects created by the AbstractFactory
with a specific Senzing configuration JSON document identifier.

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
			return wraperror.Errorf(err, "szDiagnostic.Reinitialize(%d)", configID)
		}
	}

	if factory.isSzengineInitialized {
		szEngine := &szengine.Szengine{}

		err = szEngine.Reinitialize(ctx, configID)
		if err != nil {
			return wraperror.Errorf(err, "szEngine.Reinitialize(%d)", configID)
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (factory *Szabstractfactory) destroySzConfigmanager(ctx context.Context) error {
	var err error

	if factory.isSzconfigmanagerInitialized {
		szConfigmanager := &szconfigmanager.Szconfigmanager{}

		err = szConfigmanager.Destroy(ctx)
		if err != nil {
			return wraperror.Errorf(err, "Destroy")
		}

		factory.isSzconfigmanagerInitialized = false
	}

	return nil
}

func (factory *Szabstractfactory) destroySzDiagnostic(ctx context.Context) error {
	var err error

	if factory.isSzdiagnosticInitialized {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}

		err = szDiagnostic.Destroy(ctx)
		if err != nil {
			return wraperror.Errorf(err, "Destroy")
		}

		factory.isSzdiagnosticInitialized = false
	}

	return nil
}

func (factory *Szabstractfactory) destroySzEngine(ctx context.Context) error {
	var err error

	if factory.isSzengineInitialized {
		szEngine := &szengine.Szengine{}

		err = szEngine.Destroy(ctx)
		if err != nil {
			return wraperror.Errorf(err, "Destroy")
		}

		factory.isSzengineInitialized = false
	}

	return nil
}

func (factory *Szabstractfactory) destroySzProduct(ctx context.Context) error {
	var err error

	if factory.isSzproductInitialized {
		szProduct := &szproduct.Szproduct{}

		err = szProduct.Destroy(ctx)
		if err != nil {
			return wraperror.Errorf(err, "Destroy")
		}

		factory.isSzproductInitialized = false
	}

	return nil
}
