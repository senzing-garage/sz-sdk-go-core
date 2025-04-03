//go:build linux

package szengine_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-core/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

const (
	instanceName = "SzEngine Test"
)

var (
	szEngineSingleton *szengine.Szengine
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

// TODO: See if there's any way to use currying to simplify the captureStdout* methods

func captureStdout(functionName func() error) (string, error) { //nolint
	// Reference:
	// https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang
	// Switch STDOUT.
	originalStdout, err := syscall.Dup(syscall.Stdout)
	handleErrorWithPanic(err)

	readFile, writeFile, _ := os.Pipe()
	fileDescriptor := int(writeFile.Fd())
	err = syscall.Dup2(fileDescriptor, syscall.Stdout)
	handleErrorWithPanic(err)

	// Call function.

	resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()

	err = syscall.Dup2(originalStdout, syscall.Stdout)

	handleErrorWithPanic(err)
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)

	return string(stdoutBuffer), resultErr
}

func captureStdoutReturningInt64(functionName func() (int64, error)) (string, int64, error) {
	// Reference:
	// https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang
	// Switch STDOUT.
	originalStdout, err := syscall.Dup(syscall.Stdout)
	handleErrorWithPanic(err)

	readFile, writeFile, _ := os.Pipe()
	fileDescriptor := int(writeFile.Fd())
	err = syscall.Dup2(fileDescriptor, syscall.Stdout)
	handleErrorWithPanic(err)

	// Call function.

	result, resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()

	err = syscall.Dup2(originalStdout, syscall.Stdout)
	handleErrorWithPanic(err)
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)

	return string(stdoutBuffer), result, resultErr
}

func captureStdoutReturningString(functionName func() (string, error)) (string, string, error) {
	// Reference:
	// https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang
	// Switch STDOUT.
	originalStdout, err := syscall.Dup(syscall.Stdout)
	handleErrorWithPanic(err)

	readFile, writeFile, _ := os.Pipe()
	fileDescriptor := int(writeFile.Fd())
	err = syscall.Dup2(fileDescriptor, syscall.Stdout)
	handleErrorWithPanic(err)

	// Call function.

	result, resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()

	err = syscall.Dup2(originalStdout, syscall.Stdout)
	handleErrorWithPanic(err)
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)

	return string(stdoutBuffer), result, resultErr
}

// func captureStdoutReturningUintptr(functionName func() (uintptr, error)) (string, uintptr, error) {
// Reference:
// https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang

// 	// Switch STDOUT.

// 	originalStdout, err := syscall.Dup(syscall.Stdout)
//  handleErrorWithPanic(err)
// 	readFile, writeFile, _ := os.Pipe()
// 	fileDescriptor := int(writeFile.Fd()) //nolint:gosec
// 	err = syscall.Dup2(fileDescriptor, syscall.Stdout)
//  handleErrorWithPanic(err)

// 	// Call function.

// 	result, resultErr := functionName()

// 	// Restore STDOUT.

// 	writeFile.Close()
// 	err = syscall.Dup2(originalStdout, syscall.Stdout)
//  handleErrorWithPanic(err)
// 	syscall.Close(originalStdout)

// 	// Return results.

// 	stdoutBuffer, _ := io.ReadAll(readFile)
// 	return string(stdoutBuffer), result, resultErr
// }

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() string {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	handleErrorWithPanic(err)

	databaseURL := "sqlite3://na:na@nowhere/" + dbTargetPath

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	handleErrorWithPanic(err)

	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szengine")
}

func getSzEngine(ctx context.Context) *szengine.Szengine {
	_ = ctx

	if szEngineSingleton == nil {
		settings := getSettings()
		szEngine := &szengine.Szengine{}
		_, err := captureStdout(func() error {
			return szEngine.Initialize(
				ctx,
				instanceName,
				settings,
				senzing.SzInitializeWithDefaultConfiguration,
				senzing.SzVerboseLogging,
			)
		})
		handleErrorWithPanic(err)

		szEngineSingleton = &szengine.Szengine{}
	}

	return szEngineSingleton
}

// func getVerboseSzEngineAsInterface(ctx context.Context) senzing.SzEngine {
// 	return getVerboseSzEngine(ctx)
// }

func getVerboseTestObject(t *testing.T) senzing.SzEngine {
	t.Helper()
	ctx := t.Context()

	return getSzEngine(ctx)
}

func handleErrorWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func setup() {
	setupDirectories()
	setupDatabase()

	err := setupSenzingConfiguration()
	handleErrorWithPanic(err)
}

func setupDatabase() {
	testDirectoryPath := getTestDirectoryPath()
	_, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	handleErrorWithPanic(err)
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	handleErrorWithPanic(err)

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	handleErrorWithPanic(err)
}

func setupDirectories() {
	testDirectoryPath := getTestDirectoryPath()
	err := os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	handleErrorWithPanic(err)
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	handleErrorWithPanic(err)
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()
	settings := getSettings()

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	_, err := captureStdout(func() error {
		return szConfig.Initialize(ctx, instanceName, settings, senzing.SzNoLogging)
	})
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfig.Destroy(ctx)) }()

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	_, err = captureStdout(func() error {
		return szConfigManager.Initialize(ctx, instanceName, settings, senzing.SzNoLogging)
	})
	handleErrorWithPanic(err)

	defer func() { handleErrorWithPanic(szConfigManager.Destroy(ctx)) }()

	// Create a Senzing configuration.

	_, err = captureStdout(func() error {
		return szConfig.ImportTemplate(ctx)
	})
	handleErrorWithPanic(err)

	// Add data sources to template Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, _, err := captureStdoutReturningString(func() (string, error) {
			return szConfig.AddDataSource(ctx, dataSourceCode)
		})
		handleErrorWithPanic(err)
	}

	// Create a string representation of the Senzing configuration.

	_, configDefinition, err := captureStdoutReturningString(func() (string, error) {
		return szConfig.Export(ctx)
	})
	handleErrorWithPanic(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	_, _, err = captureStdoutReturningInt64(func() (int64, error) {
		return szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	})
	handleErrorWithPanic(err)

	return nil
}

func teardown() {
	ctx := context.TODO()
	teardownSzEngine(ctx)
}

func teardownSzEngine(ctx context.Context) {
	_, err := captureStdout(func() error {
		return szEngineSingleton.Destroy(ctx)
	})
	handleErrorWithPanic(err)

	szEngineSingleton = nil
}
