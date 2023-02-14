# g2-sdk-go-base

## :warning: WARNING: g2-sdk-go-base is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing g2-sdk-go-base packages provide a
[Go](https://go.dev/)
language Software Development Kit that wraps the
Senzing C SDK APIs.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/g2-sdk-go-base.svg)](https://pkg.go.dev/github.com/senzing/g2-sdk-go-base)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/g2-sdk-go-base)](https://goreportcard.com/report/github.com/senzing/g2-sdk-go-base)
[![go-test.yaml](https://github.com/Senzing/g2-sdk-go-base/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/g2-sdk-go-base/actions/workflows/go-test.yaml)

## Overview

The Senzing g2-sdk-go-base packages enable Go programs to call Senzing library functions.
Under the covers, Golang's CGO is used by the g2-sdk-go-base packages to make calls
to the functions in the Senzing C libraries.
The `g2-sdk-go-base` implementation of the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interface is used to call the Senzing C SDK APIs directly using Go's CGO.

Other implementations of the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interface include:

- [g2-sdk-go-mock](https://github.com/Senzing/g2-sdk-go-mock) - for
  unit testing calls to the Senzing Go SDK
- [g2-sdk-go-grpc](https://github.com/Senzing/g2-sdk-go-grpc) - for
  calling Senzing SDK APIs over [gRPC](https://grpc.io/)
- [go-sdk-abstract-factory](https://github.com/Senzing/go-sdk-abstract-factory) - An
  [abstract factory pattern](https://en.wikipedia.org/wiki/Abstract_factory_pattern)
  for switching among implementations

## Developing with g2-sdk-go-base

### Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=g2-sdk-go-base
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in
   [clone-repository](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

### Install go

1. See Go's [Download and install](https://go.dev/doc/install)

### Install Senzing library

Since the Senzing library is a prerequisite, it must be installed first.
This can be done by installing the Senzing package using `apt`, `yum`,
or a technique using Docker containers.
Once complete, the Senzing library will be installed in the `/opt/senzing` directory.
This is important as the compiling of the code expects Senzing to be in `/opt/senzing`.

- Using `apt`:

    ```console
    wget https://senzing-production-apt.s3.amazonaws.com/senzingrepo_1.0.0-1_amd64.deb
    sudo apt install ./senzingrepo_1.0.0-1_amd64.deb
    sudo apt update
    sudo apt install senzingapi

    ```

- Using `yum`:

    ```console
    sudo yum install https://senzing-production-yum.s3.amazonaws.com/senzingrepo-1.0.0-1.x86_64.rpm
    sudo yum install senzingapi

    ```

- Using Docker:

  This technique can be handy if you are using MacOS or Windows and cross-compiling.

    1. Build Senzing installer.

        ```console
        curl -X GET \
            --output /tmp/senzing-versions-stable.sh \
            https://raw.githubusercontent.com/Senzing/knowledge-base/main/lists/senzing-versions-stable.sh
        source /tmp/senzing-versions-stable.sh

        sudo docker build \
            --build-arg SENZING_ACCEPT_EULA=I_ACCEPT_THE_SENZING_EULA \
            --build-arg SENZING_APT_INSTALL_PACKAGE=senzingapi=${SENZING_VERSION_SENZINGAPI_BUILD} \
            --build-arg SENZING_DATA_VERSION=${SENZING_VERSION_SENZINGDATA} \
            --no-cache \
            --tag senzing/installer:${SENZING_VERSION_SENZINGAPI} \
            https://github.com/senzing/docker-installer.git#main

        ```

    1. Install Senzing.

        ```console
            curl -X GET \
                --output /tmp/senzing-versions-stable.sh \
                https://raw.githubusercontent.com/Senzing/knowledge-base/main/lists/senzing-versions-stable.sh
            source /tmp/senzing-versions-stable.sh

            sudo rm -rf /opt/senzing
            sudo mkdir -p /opt/senzing

            sudo docker run \
                --rm \
                --user 0 \
                --volume /opt/senzing:/opt/senzing \
                senzing/installer:${SENZING_VERSION_SENZINGAPI}

        ```

### Configure Senzing

1. Move the "versioned" Senzing data to the system location.
   Example:

    ```console
      sudo mv /opt/senzing/data/3.0.0/* /opt/senzing/data/

    ```

1. Create initial configuration.
   Example:

    ```console
      sudo mkdir /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/cfgVariant.json     /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/customGn.txt        /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/customOn.txt        /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/customSn.txt        /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/defaultGNRCP.config /etc/opt/senzing
      sudo cp /opt/senzing/g2/resources/templates/stb.config          /etc/opt/senzing

    ```

### Test using SQLite database

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```

1. **Optional:** View the SQLite database.
   Example:

    ```console
    docker run \
        --env SQLITE_DATABASE=G2C.db \
        --interactive \
        --publish 9174:8080 \
        --rm \
        --tty \
        --volume /tmp/sqlite:/data \
        coleifer/sqlite-web

    ```

   Visit [localhost:9174](http://localhost:9174).

### Test using Docker-compose stack with PostgreSql database

The following instructions show how to bring up a test stack to be used
in testing the `g2-sdk-go-base` packages.

1. Identify a directory to place docker-compose artifacts.
   The directory specified will be deleted and re-created.
   Example:

    ```console
    export SENZING_DEMO_DIR=~/my-senzing-demo

    ```

1. Bring up the docker-compose stack.
   Example:

    ```console
    export PGADMIN_DIR=${SENZING_DEMO_DIR}/pgadmin
    export POSTGRES_DIR=${SENZING_DEMO_DIR}/postgres
    export RABBITMQ_DIR=${SENZING_DEMO_DIR}/rabbitmq
    export SENZING_VAR_DIR=${SENZING_DEMO_DIR}/var
    export SENZING_UID=$(id -u)
    export SENZING_GID=$(id -g)

    rm -rf ${SENZING_DEMO_DIR:-/tmp/nowhere/for/safety}
    mkdir ${SENZING_DEMO_DIR}
    mkdir -p ${PGADMIN_DIR} ${POSTGRES_DIR} ${RABBITMQ_DIR} ${SENZING_VAR_DIR}
    chmod -R 777 ${SENZING_DEMO_DIR}

    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-versions-stable.sh \
        <https://raw.githubusercontent.com/Senzing/knowledge-base/main/lists/docker-versions-stable.sh>
    source ${SENZING_DEMO_DIR}/docker-versions-stable.sh
    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-compose.yaml \
        "https://raw.githubusercontent.com/Senzing/docker-compose-demo/main/resources/postgresql/docker-compose-postgresql.yaml"

    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose up

    ```

1. In a separate terminal window, set environment variables.
   Identify Database URL of database in docker-compose stack.
   Example:

    ```console
    export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/Senzing/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2

    ```

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```

1. **Optional:** View the PostgreSQL database.
   Visit [localhost:9171](http://localhost:9171).

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
