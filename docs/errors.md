# sz-sdk-go-core errors

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6001` - szconfig
1. `6002` - szconfigmanager
1. `6003` - szdiagnostic
1. `6004` - szengine
1. `6005` - szhasher
1. `6006` - szproduct
1. `6007` - sssadm

## Common errors

### Postgresql

1. "Error: pq: SSL is not enabled on the server"
    1. The database URL needs the `sslmode` parameter.
       Example:

        ```console
        postgresql://username:password@postgres.example.com:5432/G2/?sslmode=disable
        ```

    1. [Connection String Parameters](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)

## Errors by ID

### senzing-60010001

- Trace the entering of `szconfig.AddDataSource`
- See `AddDataSource` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010002

- Trace the exiting of `szconfig.AddDataSource`
- See `AddDataSource` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010003

- Trace the entering of `szconfig.ClearLastException`
- See `ClearLastException` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010004

- Trace the exiting of `szconfig.ClearLastException`
- See `ClearLastException` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010005

- Trace the entering of `szconfig.Close`
- See `Close` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010006

- Trace the exiting of `szconfig.Close`
- See `Close` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010007

- Trace the entering of `szconfig.Create`
- See `Create` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010008

- Trace the exiting of `szconfig.Create`
- See `Create` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010009

- Trace the entering of `szconfig.DeleteDataSource`
- See `DeleteDataSource` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010010

- Trace the exiting of `szconfig.DeleteDataSource`
- See `DeleteDataSource` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010011

- Trace the entering of `szconfig.Destroy`
- See `Destroy` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010012

- Trace the exiting of `szconfig.Destroy`
- See `Destroy` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010013

- Trace the entering of `szconfig.GetLastException`
- See `GetLastException` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010014

- Trace the exiting of `szconfig.GetLastException`
- See `GetLastException` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010015

- Trace the entering of `szconfig.GetLastExceptionCode`
- See `GetLastExceptionCode` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010016

- Trace the exiting of `szconfig.GetLastExceptionCode`
- See `GetLastExceptionCode` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010017

- Trace the entering of `szconfig.Init`
- See `Init` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010018

- Trace the exiting of `szconfig.Init`
- See `Init` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010019

- Trace the entering of `szconfig.ListDataSources`
- See `ListDataSources` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010020

- Trace the exiting of `szconfig.ListDataSources`
- See `ListDataSources` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010021

- Trace the entering of `szconfig.Load`
- See `` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010022

- Trace the exiting of `szconfig.Load`
- See `Load` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010023

- Trace the entering of `szconfig.Save`
- See `Load` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010024

- Trace the exiting of `szconfig.Save`
- See `Save` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010025

- Trace the entering of `szconfig.SetLogLevel`
- See `SetLogLevel` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010026

- Trace the exiting of `szconfig.SetLogLevel`
- See `SetLogLevel` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010027

- Trace the entering of `szconfig.RegisterObserver`
- See `RegisterObserver` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010028

- Trace the exiting of `szconfig.RegisterObserver`
- See `RegisterObserver` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010029

- Trace the entering of `szconfig.UnregisterObserver`
- See `UnregisterObserver` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010030

- Trace the exiting of `szconfig.UnregisterObserver`
- See `UnregisterObserver` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010031

- Trace the entering of `szconfig.GetSdkId`
- See `GetSdkId` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60010032

- Trace the exiting of `szconfig.GetSdkId`
- See `GetSdkId` in <https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig.go>

### senzing-60014001

- Call to `G2Config_addDataSource()` failed.

### senzing-60014002

- Call to `G2Config_close()` failed.

### senzing-60014003

- Call to `G2Config_create()` failed.

### senzing-60014004

- Call to `G2Config_deleteDataSource()` failed.

### senzing-60014005

- Call to `G2Config_getLastException()` failed.

### senzing-60014006

- Call to `G2Config_destroy()` failed.

### senzing-60014007

- Call to `G2Config_init()` failed.

### senzing-60014008

- Call to `G2Config_listDataSources()` failed.

### senzing-60014009

- Call to `G2Config_load()` failed.

### senzing-60014010

- Call to `G2Config_save()` failed.
