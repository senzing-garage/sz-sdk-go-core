## g2-sdk-go-base errors

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6001` - g2config
1. `6002` - g2configmgr
1. `6003` - g2diagnostic
1. `6004` - g2engine
1. `6005` - g2hasher
1. `6006` - g2product
1. `6007` - g2ssadm

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

- Trace the entering of `AddDataSource`
- See `AddDataSource` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010002

- Trace the exiting of `AddDataSource`
- See `AddDataSource` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010003

- Trace the entering of `ClearLastException`
- See `ClearLastException` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010004

- Trace the exiting of `ClearLastException`
- See `ClearLastException` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010005

- Trace the entering of `Close`
- See `Close` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010006

- Trace the exiting of `Close`
- See `Close` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010007

- Trace the entering of `Create`
- See `Create` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010008

- Trace the exiting of `Create`
- See `Create` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010009

- Trace the entering of `DeleteDataSource`
- See `DeleteDataSource` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010010

- Trace the exiting of `DeleteDataSource`
- See `DeleteDataSource` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010011

- Trace the entering of `Destroy`
- See `Destroy` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010012

- Trace the exiting of `Destroy`
- See `Destroy` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010013

- Trace the entering of `GetLastException`
- See `GetLastException` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010014

- Trace the exiting of `GetLastException`
- See `GetLastException` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010015

- Trace the entering of `GetLastExceptionCode`
- See `GetLastExceptionCode` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010016

- Trace the exiting of `GetLastExceptionCode`
- See `GetLastExceptionCode` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010017

- Trace the entering of `Init`
- See `Init` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010018

- Trace the exiting of `Init`
- See `Init` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010019

- Trace the entering of `ListDataSources`
- See `ListDataSources` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010020

- Trace the exiting of `ListDataSources`
- See `ListDataSources` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010021

- Trace the entering of `Load`
- See `` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010022

- Trace the exiting of `Load`
- See `Load` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010023

- Trace the entering of `Save`
- See `Load` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010024

- Trace the exiting of `Save`
- See `Save` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010025

- Trace the entering of `SetLogLevel`
- See `SetLogLevel` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010026

- Trace the exiting of `SetLogLevel`
- See `SetLogLevel` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010027

- Trace the entering of `RegisterObserver`
- See `RegisterObserver` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010028

- Trace the exiting of `RegisterObserver`
- See `RegisterObserver` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010029

- Trace the entering of `UnregisterObserver`
- See `UnregisterObserver` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010030

- Trace the exiting of `UnregisterObserver`
- See `UnregisterObserver` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010031

- Trace the entering of `GetSdkId`
- See `GetSdkId` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60010032

- Trace the exiting of `GetSdkId`
- See `GetSdkId` in <https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config.go>

### senzing-60014001

- `G2Config_addDataSource()` failed.

### senzing-60014002

- `G2Config_close()` failed.

### senzing-60014003

- `G2Config_create()` failed.

### senzing-60014004

- `G2Config_deleteDataSource()` failed.

### senzing-60014005

- `G2Config_getLastException()` failed.

### senzing-60014006

- `G2Config_destroy()` failed.

### senzing-60014007

- `G2Config_init()` failed.

### senzing-60014008

- `G2Config_listDataSources()` failed.

### senzing-60014009

- `G2Config_load()` failed.

### senzing-60014010

- `G2Config_save()` failed.
