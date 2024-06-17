//go:build linux

package szengine

import (
	"context"
	"errors"
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
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

const (
	instanceName = "SzEngine Test"
)

var (
	szEngineSingleton senzing.SzEngine
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

// TODO: See if there's any way to use currying to simplify the captureStdout* methods

func captureStdout(functionName func() error) (string, error) {
	// Reference: https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang

	// Switch STDOUT.

	originalStdout, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		return "", err
	}
	readFile, writeFile, _ := os.Pipe()
	err = syscall.Dup2(int(writeFile.Fd()), syscall.Stdout)
	if err != nil {
		return "", err
	}

	// Call function.

	resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()
	err = syscall.Dup2(originalStdout, syscall.Stdout)
	if err != nil {
		return "", err
	}
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)
	return string(stdoutBuffer), resultErr
}

func captureStdoutReturningInt64(functionName func() (int64, error)) (string, int64, error) {
	// Reference: https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang

	// Switch STDOUT.

	originalStdout, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		return "", 0, err
	}
	readFile, writeFile, _ := os.Pipe()
	err = syscall.Dup2(int(writeFile.Fd()), syscall.Stdout)
	if err != nil {
		return "", 0, err
	}

	// Call function.

	result, resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()
	err = syscall.Dup2(originalStdout, syscall.Stdout)
	if err != nil {
		return "", 0, err
	}
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)
	return string(stdoutBuffer), result, resultErr
}

func captureStdoutReturningString(functionName func() (string, error)) (string, string, error) {
	// Reference: https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang

	// Switch STDOUT.

	originalStdout, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		return "", "", err
	}
	readFile, writeFile, _ := os.Pipe()
	err = syscall.Dup2(int(writeFile.Fd()), syscall.Stdout)
	if err != nil {
		return "", "", err
	}

	// Call function.

	result, resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()
	err = syscall.Dup2(originalStdout, syscall.Stdout)
	if err != nil {
		return "", "", err
	}
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)
	return string(stdoutBuffer), result, resultErr
}

func captureStdoutReturningUintptr(functionName func() (uintptr, error)) (string, uintptr, error) {
	// Reference: https://stackoverflow.com/questions/76565007/how-to-capture-the-contents-of-stderr-in-a-c-function-call-from-golang

	// Switch STDOUT.

	originalStdout, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		return "", 0, err
	}
	readFile, writeFile, _ := os.Pipe()
	err = syscall.Dup2(int(writeFile.Fd()), syscall.Stdout)
	if err != nil {
		return "", 0, err
	}

	// Call function.

	result, resultErr := functionName()

	// Restore STDOUT.

	writeFile.Close()
	err = syscall.Dup2(originalStdout, syscall.Stdout)
	if err != nil {
		return "", 0, err
	}
	syscall.Close(originalStdout)

	// Return results.

	stdoutBuffer, _ := io.ReadAll(readFile)
	return string(stdoutBuffer), result, resultErr
}

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() (string, error) {

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", err
	}
	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	settings, err := settings.BuildSimpleSettingsUsingMap(configAttrMap)
	if err != nil {
		panic(err)
	}
	return settings, err
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szengine")
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	_ = ctx
	if szEngineSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		szEngine := &szengine.Szengine{}
		_, err = captureStdout(func() error {
			return szEngine.Initialize(ctx, instanceName, settings, senzing.SzInitializeWithDefaultConfiguration, senzing.SzVerboseLogging)
		})
		if err != nil {
			fmt.Println(err)
		}
		szEngineSingleton = szEngine
	}
	return szEngineSingleton
}

// func getVerboseSzEngineAsInterface(ctx context.Context) senzing.SzEngine {
// 	return getVerboseSzEngine(ctx)
// }

func getVerboseTestObject(ctx context.Context, test *testing.T) senzing.SzEngine {
	_ = test
	return getSzEngine(ctx)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error

	err = setupDirectories()
	if err != nil {
		return fmt.Errorf("Failed to set up directories. Error: %w", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("Failed to set up database. Error: %w", err)
	}
	err = setupSenzingConfiguration()
	if err != nil {
		panic(err)
	}
	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
	}
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		return fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			databaseTemplatePath, err.Error())
	}

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	if err != nil {
		return fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
			databaseTemplatePath, testDirectoryPath, err)
	}
	return err
}

func setupDirectories() error {
	var err error
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}
	return err
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	settings, err := getSettings()
	if err != nil {
		panic(err)
	}

	// Create sz objects.

	szConfig := &szconfig.Szconfig{}
	_, err = captureStdout(func() error {
		return szConfig.Initialize(ctx, instanceName, settings, senzing.SzNoLogging)
	})
	if err != nil {
		panic(err)
	}
	defer func() { handleError(szConfig.Destroy(ctx)) }()

	// Create an in memory Senzing configuration.

	_, configHandle, err := captureStdoutReturningUintptr(func() (uintptr, error) {
		return szConfig.CreateConfig(ctx)
	})
	if err != nil {
		panic(err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, _, err := captureStdoutReturningString(func() (string, error) {
			return szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		})
		if err != nil {
			panic(err)
		}
	}

	// Create a string representation of the in-memory configuration.

	_, configDefinition, err := captureStdoutReturningString(func() (string, error) {
		return szConfig.ExportConfig(ctx, configHandle)
	})
	if err != nil {
		panic(err)
	}

	// Close szConfig in-memory object.

	_, err = captureStdout(func() error {
		return szConfig.CloseConfig(ctx, configHandle)
	})
	if err != nil {
		panic(err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	szConfigManager := &szconfigmanager.Szconfigmanager{}
	_, err = captureStdout(func() error {
		return szConfigManager.Initialize(ctx, instanceName, settings, senzing.SzNoLogging)
	})
	if err != nil {
		panic(err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	_, configID, err := captureStdoutReturningInt64(func() (int64, error) {
		return szConfigManager.AddConfig(ctx, configDefinition, configComment)
	})
	if err != nil {
		panic(err)
	}

	_, err = captureStdout(func() error {
		return szConfigManager.SetDefaultConfigID(ctx, configID)
	})
	if err != nil {
		panic(err)
	}

	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzEngine(ctx)
	return err
}

func teardownSzEngine(ctx context.Context) error {
	_, err := captureStdout(func() error {
		return szEngineSingleton.Destroy(ctx)
	})
	if err != nil {
		return err
	}
	szEngineSingleton = nil
	return nil
}
