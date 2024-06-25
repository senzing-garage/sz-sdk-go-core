# sz-sdk-go-core

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: sz-sdk-go-core is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `sz-sdk-go-core` packages provide a [Go]
language Software Development Kit that wraps the
Senzing C SDK APIs.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/sz-sdk-go-core.svg)](https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing-garage/sz-sdk-go-core)](https://goreportcard.com/report/github.com/senzing-garage/sz-sdk-go-core)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/senzing-garage/sz-sdk-go-core/blob/main/LICENSE)

[![gosec.yaml](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/gosec.yaml/badge.svg)](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/gosec.yaml)
[![go-test-linux.yaml](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-linux.yaml/badge.svg)](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-linux.yaml)
[![go-test-darwin.yaml](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-darwin.yaml/badge.svg)](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-darwin.yaml)
[![go-test-windows.yaml](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-windows.yaml/badge.svg)](https://github.com/senzing-garage/sz-sdk-go-core/actions/workflows/go-test-windows.yaml)

## Overview

The Senzing `sz-sdk-go-core` packages enable Go programs to call Senzing library functions.
Under the covers, Golang's CGO is used by the sz-sdk-go-core packages to make calls
to the functions in the Senzing C libraries.
The `sz-sdk-go-core` implementation of the [sz-sdk-go]
interface is used to call the Senzing C SDK APIs directly using Go's CGO.

Other implementations of the [sz-sdk-go]
interface include:

- [sz-sdk-go-mock] - for unit testing calls to the Senzing Go SDK
- [sz-sdk-go-grpc] - for  calling Senzing SDK APIs over [gRPC]
- [go-sdk-abstract-factory] - An [abstract factory pattern]for switching among implementations

## Use

(TODO:)

## References

1. [Development]
1. [Errors]
1. [Examples]
1. [Package reference]

[Go]: https://go.dev/
[Senzing]: https://senzing.com/
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing Garage]: https://github.com/senzing-garage-garage
[sz-sdk-go]: https://github.com/senzing-garage/sz-sdk-go
[sz-sdk-go-mock]: https://github.com/senzing-garage/sz-sdk-go-mock
[sz-sdk-go-grpc]: https://github.com/senzing-garage/sz-sdk-go-grpc
[go-sdk-abstract-factory]: https://github.com/senzing-garage/go-sdk-abstract-factory
[abstract factory pattern]: https://en.wikipedia.org/wiki/Abstract_factory_pattern
[gRPC]: https://grpc.io/
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-core
