#!/bin/bash -e
current_release_version () {
  bosh2 -e vbox releases|grep self-service-adapter|head -1|cut  -f2|cut -d "." -f2|cut -d "*" -f1
}
mkdir -p dev_releases
rm -rf ./dev_releases/*/*gz
newVersion=$((`current_release_version`+1))
echo $newVersion
v=0+dev.$newVersion

bosh2 create-release --version=$v --force --tarball=./dev_releases/self-service-adapter/self-service-adapter-release-${v}.tgz
bosh2 -e vbox upload-release ./dev_releases/self-service-adapter/self-service-adapter-release-${v}.tgz
