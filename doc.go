//nolint
/*
Module sz-sdk-go-core implements Senzing objects that communicate with the Senzing native C binary, libSz.so.

# Synopsis

The Senzing sz-sdk-go-core packages are wrappers over Senzing's C-based library.

# Overview

The Senzing sz-sdk-go-core packages enable Go programs to call Senzing library functions.
Under the covers, Golang's CGO is used by the sz-sdk-go-core packages to make the calls
to the Senzing functions.

More information at [sz-sdk-go-core].

# Installing Senzing library

Since the Senzing API library is a pre-requisite, it must be installed first.
This can be done by installing the Senzing package using apt, yum,
or a technique using Docker containers.
Once complete, the Senzing library will be installed in the /opt/senzing directory.

See [How to install Senzing API].

# Examples

In addition to the examples in the documentation, the test files provide additional examples:

  - [szconfig_test.go]
  - [szconfigmanager_test.go]
  - [szdiagnostic_test.go]
  - [szengine_test.go]
  - [szproduct_test.go]

[How to install Senzing API]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-api.md
[sz-sdk-go-core]: https://github.com/senzing-garage/sz-sdk-go-core
[szconfig_test.go]: https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_test.go
[szconfigmanager_test.go]: https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_test.go
[szdiagnostic_test.go]: https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_test.go
[szengine_test.go]: https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_test.go
[szproduct_test.go]: https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szproduct/szproduct_test.go
*/
package main
