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
	ConfigID       int64
	InstanceName   string
	isClosed       bool
	mutex          sync.Mutex
	once           sync.Once
	semaphores     []*szengine.Szengine
	Settings       string
	VerboseLogging int64
}

// ----------------------------------------------------------------------------
// senzing.SzAbstractFactory interface methods
// ----------------------------------------------------------------------------

/*
Method Close prevents the AbstractFactory from creating any more object.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Close(ctx context.Context) error {
	var err error

	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	factory.isClosed = true

	for _, semaphore := range factory.semaphores {
		err = semaphore.Destroy(ctx)
	}
	factory.semaphores = nil

	return wraperror.Errorf(err, wraperror.NoMessage)
}

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
		return result, wraperror.Errorf(errForPackage, "SzAbstractFactory is closed")
	}

	factory.once.Do(func() {
		err = factory.initializeAbstractFactory(ctx)
	})

	if err != nil {
		factory.once = sync.Once{}
		return result, wraperror.Errorf(
			err,
			"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzConfigManager]",
		)
	}

	result = &szconfigmanager.Szconfigmanager{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)

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
		return result, wraperror.Errorf(errForPackage, "SzAbstractFactory is closed")
	}

	factory.once.Do(func() {
		err = factory.initializeAbstractFactory(ctx)
	})

	if err != nil {
		factory.once = sync.Once{}

		return result, wraperror.Errorf(
			err,
			"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzDiagnostic]",
		)
	}

	result = &szdiagnostic.Szdiagnostic{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)

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
		return result, wraperror.Errorf(errForPackage, "SzAbstractFactory is closed")
	}

	factory.once.Do(func() {
		err = factory.initializeAbstractFactory(ctx)
	})

	if err != nil {
		factory.once = sync.Once{}

		return result, wraperror.Errorf(
			err,
			"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzEngine]",
		)
	}

	result = &szengine.Szengine{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.ConfigID, factory.VerboseLogging)

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
		return result, wraperror.Errorf(errForPackage, "SzAbstractFactory is closed")
	}

	factory.once.Do(func() {
		err = factory.initializeAbstractFactory(ctx)
	})

	if err != nil {
		factory.once = sync.Once{}

		return result, wraperror.Errorf(
			err,
			"Cannot create AbstractFactory until prior AbstractFactory has been closed and objects created by that factory destroyed [SzProduct]",
		)
	}

	result = &szproduct.Szproduct{}
	err = result.Initialize(ctx, factory.InstanceName, factory.Settings, factory.VerboseLogging)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
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

	if factory.isClosed {
		return wraperror.Errorf(errForPackage, "SzAbstractFactory is closed")
	}

	factory.ConfigID = configID

	if factory.szDiagnosticExists(ctx) {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}

		err = szDiagnostic.Reinitialize(ctx, configID)
		if err != nil {
			return wraperror.Errorf(err, "szDiagnostic.Reinitialize(%d)", configID)
		}
	}

	if factory.szEngineExists(ctx) {
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
Method initializeAbstractFactory performs first-time checking and ...

A returned error signifies that there are still objects registered with the Senzing engine, et al.
*/
func (factory *Szabstractfactory) initializeAbstractFactory(ctx context.Context) error {
	var err error

	err = factory.verifyNoSenzingObjects(ctx)
	if err != nil {
		return wraperror.Errorf(err, "verifyNoSenzingObjects")
	}

	// IMPROVE:  At this point in the code, there is a slight concurrency hole.
	// If multiple AbstractFactories pass the verifyNoSenzingObjects test,
	// it's a race condition to see which configuration wins.
	// This may mitigated by use of runtime.LockOSThread()

	// Create semaphore.

	semaphoreSzEngine := &szengine.Szengine{}
	err = semaphoreSzEngine.Initialize(
		ctx,
		factory.InstanceName,
		factory.Settings,
		factory.ConfigID,
		factory.VerboseLogging,
	)

	if factory.semaphores == nil {
		factory.semaphores = []*szengine.Szengine{}
	}
	factory.semaphores = append(factory.semaphores, semaphoreSzEngine)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (factory *Szabstractfactory) szConfigManagerExists(ctx context.Context) bool {
	_ = ctx
	szConfigManager := &szconfigmanager.Szconfigmanager{}
	_, err := szConfigManager.GetDefaultConfigID(ctx)

	return err == nil
}

func (factory *Szabstractfactory) szDiagnosticExists(ctx context.Context) bool {
	_ = ctx
	szDiagnostic := &szdiagnostic.Szdiagnostic{}
	_, err := szDiagnostic.GetRepositoryInfo(ctx)

	return err == nil
}

func (factory *Szabstractfactory) szEngineExists(ctx context.Context) bool {
	_ = ctx
	szEngine := &szengine.Szengine{}
	_, err := szEngine.GetActiveConfigID(ctx)

	return err == nil
}

func (factory *Szabstractfactory) szProductExists(ctx context.Context) bool {
	_ = ctx

	// IMPROVE: Is there a way to check for the existence.

	return false
}

/*
Method verifyNoSenzingObjects determines if any Senzing objects are registered with the
underlying C binaries.

A returned error signifies that there are still objects registered with the Senzing engine, et al.
*/
func (factory *Szabstractfactory) verifyNoSenzingObjects(ctx context.Context) error {
	var err error

	if factory.szConfigManagerExists(ctx) {
		return wraperror.Errorf(
			errForPackage,
			"Must call Destroy() on existing SzConfigManager instances before creating new SzAbstractFactory.",
		)
	}

	if factory.szDiagnosticExists(ctx) {
		return wraperror.Errorf(
			errForPackage,
			"Must call Destroy() on existing SzDiagnostic instances before creating new SzAbstractFactory.",
		)
	}

	if factory.szEngineExists(ctx) {
		return wraperror.Errorf(
			errForPackage,
			"Must call Destroy() on existing SzEngine instances before creating new SzAbstractFactory.",
		)
	}

	if factory.szProductExists(ctx) {
		return wraperror.Errorf(
			errForPackage,
			"Must call Destroy() on existing SzProduct instances before creating new SzAbstractFactory.",
		)
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}
