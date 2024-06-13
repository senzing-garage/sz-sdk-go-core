//go:build linux

package szengine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"testing"

	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
)

var (
	szEngineVerboseSingleton *Szengine
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzengine_AddRecord_parameterCheck(test *testing.T) {

	// file, err := os.Open("/dev/stdout") // For read access.
	// require.NoError(test, err)
	// defer file.Close()

	// output := captureOutput(func() {
	// 	ctx := context.TODO()
	// 	szEngine := getVerboseTestObject(ctx, test)
	// 	flags := senzing.SzWithoutInfo
	// 	records := []record.Record{
	// 		truthset.CustomerRecords["1001"],
	// 		truthset.CustomerRecords["1002"],
	// 	}
	// 	for _, record := range records {
	// 		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
	// 		require.NoError(test, err)
	// 		printActual(test, actual)
	// 	}
	// 	for _, record := range records {
	// 		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	// 		require.NoError(test, err)
	// 		printActual(test, actual)
	// 	}
	// })

	// use syscall.Dup to get a copy of stderr
	origStderr, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		panic(err)
	}

	r, w, _ := os.Pipe()

	// Clone the pipe's writer to the actual Stderr descriptor; from this point
	// on, writes to Stderr will go to w.
	if err = syscall.Dup2(int(w.Fd()), syscall.Stdout); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	szEngine := getVerboseTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}

	fmt.Printf("\n\n\n\n\n\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n\n")

	w.Close()
	syscall.Dup2(origStderr, syscall.Stdout)
	syscall.Close(origStderr)

	b, _ := io.ReadAll(r)

	fmt.Printf("\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> got output: %s\n", string(b))
	fmt.Fprintf(os.Stderr, "\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>2> got output: %s\n", string(b))
	fmt.Fprintf(os.Stderr, "stderr works normally\n")

	// aBuffer := make([]byte, 5000)
	// n1, err := file.Read(aBuffer)
	// require.NoError(test, err)

	// fmt.Printf("\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>> %d bytes: %s\n", n1, string(aBuffer[:n1]))

	// fmt.Printf("\n>>>>>>>>>>> output: %s\n>>>>>>>>>>>>\n", output)

}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func InterceptStdout() (*os.File, *os.File, func()) {
	backupStd := os.Stdout
	backupErr := os.Stderr
	r, w, _ := os.Pipe()
	//Restore streams
	cleanup := func() {
		os.Stdout = backupStd
		os.Stderr = backupErr
	}
	os.Stdout = w
	os.Stderr = w
	return r, w, cleanup
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func getVerboseSzEngine(ctx context.Context) *Szengine {
	_ = ctx
	if szEngineVerboseSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			fmt.Printf("getSettings() Error: %v\n", err)
			return nil
		}
		szEngineVerboseSingleton = &Szengine{}
		err = szEngineVerboseSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			fmt.Printf("SetLogLevel() Error: %v\n", err)
			return nil
		}
		if logLevel == "TRACE" {
			szEngineVerboseSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineVerboseSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				fmt.Printf("RegisterObserver() Error: %v\n", err)
				return nil
			}
			err = szEngineVerboseSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				fmt.Printf("SetLogLevel() - 2 Error: %v\n", err)
				return nil
			}
		}
		err = szEngineVerboseSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), senzing.SzVerboseLogging)
		if err != nil {
			fmt.Println(err)
		}
	}
	return szEngineVerboseSingleton
}

func getVerboseSzEngineAsInterface(ctx context.Context) senzing.SzEngine {
	return getVerboseSzEngine(ctx)
}

func getVerboseTestObject(ctx context.Context, test *testing.T) *Szengine {
	_ = test
	return getVerboseSzEngine(ctx)
}
