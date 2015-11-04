# MAKEFILE
#
# @author      Nicola Asuni <nicola.asuni@datasift.com>
# @link        https://github.com/datasift/GoOse
# ------------------------------------------------------------------------------

# List special make targets that are not associated with files
.PHONY: help all test format fmtcheck vet lint coverage qa deps install uninstall clean nuke build rpm dist

# Ensure everyone is using bash. Note that Ubuntu now uses dash which doesn't support PIPESTATUS.
SHELL=/bin/bash

# name of RPM or DEB package
PKGNAME=GoOse

# Go lang path
GOPATH=$(shell readlink -f $(shell pwd)/../../../../)

# Message to diplay for commands that are not package-compatible
NOPKGCMD="This command is not supported as this is a Go package and not a stand-alone application."

# --- MAKE TARGETS ---

# Display general help about this command
help:
	@echo ""
	@echo "Welcome to $(PKGNAME) make."
	@echo "The following commands are available:"
	@echo ""
	@echo "    make qa         : Run all the tests"
	@echo ""
	@echo "    make test       : Run the unit tests"
	@echo "    make test.short : Run the unit tests with the short option"
	@echo ""
	@echo "    make format     : Format the source code"
	@echo "    make fmtcheck   : Check if the source code has been formatted"
	@echo "    make vet        : Check for syntax errors"
	@echo "    make lint       : Check for style errors"
	@echo "    make coverage   : Generate the coverage report"
	@echo ""
	@echo "    make docs       : Generate source code documentation"
	@echo ""
	@echo "    make deps       : Get the dependencies"
	@echo "    make nuke       : Deletes any intermediate file"
	@echo ""

# Alias for help target
all: help

# Run the unit tests
test:
	@mkdir -p target/test
	GOPATH=$(GOPATH) go test -race -v ./... | tee >(PATH=$(GOPATH)/bin:$(PATH) go-junit-report > target/test/report.xml); test $${PIPESTATUS[0]} -eq 0

# Run the unit tests with the short option
test.short:
	@mkdir -p target/test
	GOPATH=$(GOPATH) go test -short -race -v ./... | tee >(PATH=$(GOPATH)/bin:$(PATH) go-junit-report > target/test/report.xml); test $${PIPESTATUS[0]} -eq 0

# Format the source code
format:
	@find ./ -type f -name "*.go" -exec gofmt -w {} \;

# Check if the source code has been formatted
fmtcheck:
	@mkdir -p target
	@find ./ -type f -name "*.go" -exec gofmt -d {} \; | tee target/format.diff
	@test ! -s target/format.diff || { echo "ERROR: the source code has not been formatted - please use 'make format' or 'gofmt'"; exit 1; }

# Check for syntax errors
vet:
	GOPATH=$(GOPATH) go vet ./...

# Check for style errors
lint:
	GOPATH=$(GOPATH) PATH=$(GOPATH)/bin:$(PATH) golint ./...

# Generate the coverage report
coverage:
	@mkdir -p target/report
	GOPATH=$(GOPATH) ./coverage.sh

# Generate source docs
docs:
	@mkdir -p target/docs
	nohup sh -c 'GOPATH=$(GOPATH) godoc -http=127.0.0.1:6060' > target/godoc_server.log 2>&1 &
	wget --directory-prefix=target/docs/ --execute robots=off --retry-connrefused --recursive --no-parent --adjust-extension --page-requisites --convert-links http://127.0.0.1:6060/pkg/github.com/datasift/'${PKGNAME}'/ ; kill -9 `lsof -ti :6060`
	echo '<html><head><meta http-equiv="refresh" content="0;./127.0.0.1:6060/pkg/github.com/datasift/'${PKGNAME}'/index.html"/></head><a href="./127.0.0.1:6060/pkg/github.com/datasift/'${PKGNAME}'/index.html">'${PKGNAME}' Documentation ...</a></html>' > target/docs/index.html

# Alias to run targets: fmtcheck test vet lint coverage
qa: fmtcheck test vet lint coverage

# --- INSTALL ---

# Get the dependencies
deps:
	GOPATH=$(GOPATH) go get ./...
	GOPATH=$(GOPATH) go get github.com/golang/lint/golint
	GOPATH=$(GOPATH) go get github.com/jstemmer/go-junit-report

# Install this application
install:
	@echo $(NOPKGCMD)

# Remove all installed files (excluding configuration files)
uninstall:
	@echo $(NOPKGCMD)

# Remove any build artifact
clean:
	@echo $(NOPKGCMD)

# Deletes any intermediate file
nuke:
	rm -rf ./target

# Compile the application
build:
	@echo $(NOPKGCMD)

# --- PACKAGING ---

# Build the RPM package for RedHat-like Linux distributions
rpm:
	@echo $(NOPKGCMD)

# Execute all tests and build the RPM package
dist:
	@echo $(NOPKGCMD)
