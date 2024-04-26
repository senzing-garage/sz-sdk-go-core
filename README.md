# sz-sdk-go-core

If you are beginning your journey with
[Senzing](https://senzing.com/),
please start with
[Senzing Quick Start guides](https://docs.senzing.com/quickstart/).

You are in the
[Senzing Garage](https://github.com/senzing-garage-garage)
where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: sz-sdk-go-core is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing sz-sdk-go-core packages provide a
[Go](https://go.dev/)
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
The `sz-sdk-go-core` implementation of the
[sz-sdk-go](https://github.com/senzing-garage/sz-sdk-go)
interface is used to call the Senzing C SDK APIs directly using Go's CGO.

Other implementations of the
[sz-sdk-go](https://github.com/senzing-garage/sz-sdk-go)
interface include:

- [sz-sdk-go-mock](https://github.com/senzing-garage/sz-sdk-go-mock) - for
  unit testing calls to the Senzing Go SDK
- [sz-sdk-go-grpc](https://github.com/senzing-garage/sz-sdk-go-grpc) - for
  calling Senzing SDK APIs over [gRPC](https://grpc.io/)
- [go-sdk-abstract-factory](https://github.com/senzing-garage/go-sdk-abstract-factory) - An
  [abstract factory pattern](https://en.wikipedia.org/wiki/Abstract_factory_pattern)
  for switching among implementations

## Use

(TODO:)

## References

1. [Development](docs/development.md)
1. [Errors](docs/errors.md)
1. [Examples](docs/examples.md)
1. [Package reference](https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go-core)
