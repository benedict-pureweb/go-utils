.PHONY: FORCE

date=$(shell date +'%Y%m%d')
# Get commit hash
revision=$(shell git rev-parse --short HEAD)
# Check if we are in the master branch, if not add .dev to the revision
dev=$(shell git rev-parse --abbrev-ref HEAD | grep master > /dev/null || echo .dev)
build_dir=output

package=yaml-parse

deb: ../bin/${package}
	$(eval version := $(shell ../bin/${package} --version | perl -pe 's/Version: //'))
	@echo version: ${version}
	# Create deb dir structure
	mkdir -p ${build_dir}/${package}/DEBIAN
	mkdir -p ${build_dir}/${package}/usr/bin
	# Copy debian files
	cp ./debian-control ${build_dir}/${package}/DEBIAN/control
	cp ./debian-postinst ${build_dir}/${package}/DEBIAN/postinst
	chmod 755 ${build_dir}/${package}/DEBIAN/postinst
	# Update control file
	perl -i -pe 's/VERSION/${version}/' ${build_dir}/${package}/DEBIAN/control
	perl -i -pe 's/PACKAGE/${package}/' ${build_dir}/${package}/DEBIAN/control
	perl -i -pe 's/DESCRIPTION/${package}/' ${build_dir}/${package}/DEBIAN/control
	# Copy main binary
	cp ../bin/${package} ${build_dir}/${package}/usr/bin
	cd ${build_dir} && dpkg-deb --build ${package} ${package}_${version}_amd64.deb
	@echo Output: ${build_dir}/${package}_${version}_amd64.deb

../bin/$(package): FORCE
	# Clean up
	-rm ${PWD}/../bin/${package}
	# Build main binary
	go build -v -ldflags="-X main.BuildMetadata=${date}.${revision}${dev}" -o ../bin/${package} ../cmd/${package}/main.go

FORCE:
