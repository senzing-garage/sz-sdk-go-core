package szabstractfactory

import (
	"context"
	"sync"

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
	ConfigID               int64
	InstanceName           string
	isClosed               bool
	mutex                  sync.Mutex
	Settings               string
	szConfigManagerCounter int
	szDiagnosticCounter    int
	szEngineCounter        int
	szProductCounter       int
	VerboseLogging         int64
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
	var (
		err    error
		result *szconfigmanager.Szconfigmanager
	)

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if factory.isClosed {
		return result, wraperror.Errorf(err, "SzAbstractFactory is closed")
	}

	if factory.isInitialState() {
		factory.destroy(ctx)
	}

	result = &szconfigmanager.Szconfigmanager{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	factory.szConfigManagerCounter++

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
	var (
		err    error
		result *szdiagnostic.Szdiagnostic
	)

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if factory.isClosed {
		return result, wraperror.Errorf(err, "SzAbstractFactory is closed")
	}

	if factory.isInitialState() {
		factory.destroy(ctx)
	}

	result = &szdiagnostic.Szdiagnostic{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	factory.szDiagnosticCounter++

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
	var (
		err    error
		result *szengine.Szengine
	)

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if factory.isClosed {
		return result, wraperror.Errorf(err, "SzAbstractFactory is closed")
	}

	if factory.isInitialState() {
		factory.destroy(ctx)
	}

	result = &szengine.Szengine{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)
	factory.szEngineCounter++

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
	var (
		err    error
		result *szproduct.Szproduct
	)

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if factory.isClosed {
		return result, wraperror.Errorf(err, "SzAbstractFactory is closed")
	}

	if factory.isInitialState() {
		factory.destroy(ctx)
	}

	result = &szproduct.Szproduct{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)
	factory.szProductCounter++

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Close prevents factory from creating objects.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Close(ctx context.Context) error {
	var err error

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	factory.isClosed = true

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Destroy(ctx context.Context) error {
	var err error

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	factory.destroy(ctx)

	return wraperror.Errorf(err, wraperror.NoMessage)
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

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	factory.ConfigID = configID

	if factory.szDiagnosticCounter > 0 {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}

		err = szDiagnostic.Reinitialize(ctx, configID)
		if err != nil {
			return wraperror.Errorf(err, "szDiagnostic.Reinitialize(%d)", configID)
		}
	}

	if factory.szEngineCounter > 0 {
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

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) destroy(ctx context.Context) {
	factory.destroySzConfigmanager(ctx)
	factory.destroySzDiagnostic(ctx)
	factory.destroySzEngine(ctx)
	factory.destroySzProduct(ctx)
}

func (factory *Szabstractfactory) destroySzConfigmanager(ctx context.Context) {
	var err error

	szConfigmanager := &szconfigmanager.Szconfigmanager{}

	for {
		err = szConfigmanager.Destroy(ctx)
		if err != nil {
			break
		}
	}

	factory.szConfigManagerCounter = 0
}

func (factory *Szabstractfactory) destroySzDiagnostic(ctx context.Context) {
	var err error

	szDiagnostic := &szdiagnostic.Szdiagnostic{}

	for {
		err = szDiagnostic.Destroy(ctx)
		if err != nil {
			break
		}
	}

	factory.szDiagnosticCounter = 0
}

func (factory *Szabstractfactory) destroySzEngine(ctx context.Context) {
	var err error

	szEngine := &szengine.Szengine{}

	for {
		err = szEngine.Destroy(ctx)
		if err != nil {
			break
		}
	}

	factory.szEngineCounter = 0
}

func (factory *Szabstractfactory) destroySzProduct(ctx context.Context) {
	var err error

	szProduct := &szproduct.Szproduct{}

	for {
		err = szProduct.Destroy(ctx)
		if err != nil {
			break
		}
	}

	factory.szProductCounter = 0
}

func (factory *Szabstractfactory) isInitialState() bool {
	total := factory.szConfigManagerCounter + factory.szDiagnosticCounter + factory.szEngineCounter + factory.szProductCounter
	result := total == 0

	return result
}
