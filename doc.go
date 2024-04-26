/*
# Synopsis

The Senzing sz-sdk-go-core packages are wrappers over Senzing's C-based library.

# Overview

The Senzing sz-sdk-go-core packages enable Go programs to call Senzing library functions.
Under the covers, Golang's CGO is used by the sz-sdk-go-core packages to make the calls
to the Senzing functions.

More information at https://github.com/senzing-garage/sz-sdk-go-core

# Installing Senzing library

Since the Senzing library is a pre-requisite, it must be installed first.
This can be done by installing the Senzing package using apt, yum,
or a technique using Docker containers.
Once complete, the Senzing library will be installed in the /opt/senzing directory.

Using apt:

	wget https://senzing-production-apt.s3.amazonaws.com/senzingrepo_1.0.0-1_amd64.deb
	sudo apt install ./senzingrepo_1.0.0-1_amd64.deb
	sudo apt update
	sudo apt install senzingapi

Using yum:

	sudo yum install https://senzing-production-yum.s3.amazonaws.com/senzingrepo-1.0.0-1.x86_64.rpm
	sudo yum install senzingapi

Using Docker, build an installer:

	curl -X GET \
	    --output /tmp/senzing-versions-latest.sh \
	    https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/lists/senzing-versions-latest.sh
	source /tmp/senzing-versions-latest.sh

	sudo docker build \
	    --build-arg SENZING_ACCEPT_EULA=I_ACCEPT_THE_SENZING_EULA \
	    --build-arg SENZING_APT_INSTALL_PACKAGE=senzingapi=${SENZING_VERSION_SENZINGAPI_BUILD} \
	    --build-arg SENZING_DATA_VERSION=${SENZING_VERSION_SENZINGDATA} \
	    --no-cache \
	    --tag senzing/installer:${SENZING_VERSION_SENZINGAPI} \
	    https://github.com/senzing-garage/docker-installer.git#main

Using Docker, install Senzing:

	sudo rm -rf /opt/senzing
	sudo mkdir -p /opt/senzing

	sudo docker run \
	    --rm \
	    --user 0 \
	    --volume /opt/senzing:/opt/senzing \
	    senzing/installer:${SENZING_VERSION_SENZINGAPI}

# Examples

Examples of use can be seen in the xxxx_test.go files.
*/
package main
