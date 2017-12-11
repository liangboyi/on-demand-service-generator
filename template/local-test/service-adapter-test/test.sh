#!/bin/bash -e
rm -rf bin
current_path=`pwd`
cp -r ../../self-service-adapter ./
export BOSH_INSTALL_TARGET=$current_path
PACKAGE_NAME=github.com/pivotal-cf/self-service-adapter
PACKAGE_DIR=${BOSH_INSTALL_TARGET}/src/${PACKAGE_NAME}

mkdir -p $(dirname $PACKAGE_DIR)

cp -a $(basename $PACKAGE_NAME)/ $PACKAGE_DIR

export GOROOT=/usr/local/go
export GOPATH=$BOSH_INSTALL_TARGET:${PACKAGE_DIR}/vendor
export PATH=$GOROOT/bin:$PATH
echo "-------------------set GO ENV-------------------"
echo "GOROOT="$GOROOT
echo "GOPATH="$GOPATH
echo "PATH="$PATH
echo "-------------------------------------------------"

echo "===== compling service-adapter binary ====="
echo "... ... ..."
go install ${PACKAGE_NAME}/cmd/service-adapter

# clean up source artifacts
rm -rf ${BOSH_INSTALL_TARGET}/src ${BOSH_INSTALL_TARGET}/pkg ${BOSH_INSTALL_TARGET}/self-service-adapter

echo "==== finished ===="
