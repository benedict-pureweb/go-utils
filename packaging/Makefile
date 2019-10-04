.PHONY: deb

date=$(shell date +'%Y%m%d')
revision=$(shell git rev-parse --short HEAD)
package=yaml-parse
build_dir=output

deb:
	# Build main binary
	go build -v -ldflags="-X main.BuildMetadata=${date}.${revision}" \
		-o ../bin/${package} ../cmd/${package}/main.go
	$(eval version := $(shell ../bin/${package} --version | perl -pe 's/Version: //'))
	@echo ${version}
	# Create deb dir structure
	mkdir -p ${build_dir}/${package}/DEBIAN
	mkdir -p ${build_dir}/${package}/usr/bin
	# Copy debian files
	cp ./debian-control ${build_dir}/${package}/DEBIAN/control
	perl -i -pe 's/VERSION/${version}/' ${build_dir}/${package}/DEBIAN/control
	perl -i -pe 's/PACKAGE/${package}/' ${build_dir}/${package}/DEBIAN/control
	perl -i -pe 's/DESCRIPTION/${package}/' ${build_dir}/${package}/DEBIAN/control
	cp ./debian-postinst ${build_dir}/${package}/DEBIAN/postinst
	chmod 755 ${build_dir}/${package}/DEBIAN/postinst
	# Copy main binary
	cp ../bin/${package} ${build_dir}/${package}/usr/bin
	cd ${build_dir} && dpkg-deb --build ${package} ${package}_${version}_amd64.deb
	@echo Output: ${build_dir}/${package}_${version}_amd64.deb
