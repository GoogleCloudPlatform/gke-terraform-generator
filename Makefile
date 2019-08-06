# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Make will use bash instead of sh
SHELL := /usr/bin/env bash

PROJECT?=must-override

.PHONY: all
all: build test

.PHONY: gazell-update
gazelle-update:
	@bazel run //:gazelle -- update

.PHONY: gazell-update-mod
gazelle-update-mod:
	@GO111MODULE=on go mod tidy
	@bazel run //:gazelle -- update-repos -from_file=go.mod
	@bazel run //:gazelle -- update

.PHONY: build
build:
	@bazel build //...

.PHONY: run
run:
	@bazel run //:gke-tf

.PHONY: cross-build
cross-build:
	@rm -rf dist
	@mkdir -p dist/{darwin-amd64,linux-amd64,windows-amd64}
	@bazel build --features=pure --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64 //...
	@cp bazel-bin/darwin_amd64_pure_stripped/gke-tf dist/darwin-amd64/
	@bazel build --features=pure --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //...
	@cp bazel-bin/linux_amd64_pure_stripped/gke-tf dist/linux-amd64/
	@bazel build --features=pure --platforms=@io_bazel_rules_go//go/toolchain:windows_amd64 //...
	@cp bazel-bin/windows_amd64_pure_stripped/gke-tf.exe dist/windows-amd64/
	@echo "binaries staged in dist directory"


.PHONY: test
test:
	@bazel test --test_verbose_timeout_warnings //...

.PHONY: clean
clean:
	@bazel clean
	@rm -rf dist

.PHONY: gen
gen: build
	bazel-bin/darwin_amd64_stripped/gke-tf gen -d /tmp/tf/ -f examples/full.yaml -o -p ${PROJECT}

lint: check_shell check_python check_golang check_terraform check_docker \
	check_base_files check_headers check_trailing_whitespace

# The .PHONY directive tells make that this isn't a real target and so
# the presence of a file named 'check_shell' won't cause this target to stop
# working
.PHONY: check_shell
check_shell:
	@source test/make.sh && check_shell

.PHONY: check_python
check_python:
	@source test/make.sh && check_python

.PHONY: check_golang
check_golang:
	@source test/make.sh && golang

.PHONY: check_terraform
check_terraform:
	@source test/make.sh && check_terraform

.PHONY: check_docker
check_docker:
	@source test/make.sh && docker

.PHONY: check_base_files
check_base_files:
	@source test/make.sh && basefiles

.PHONY: check_shebangs
check_shebangs:
	@source test/make.sh && check_bash

.PHONY: check_trailing_whitespace
check_trailing_whitespace:
	@source test/make.sh && check_trailing_whitespace

.PHONY: check_headers
check_headers:
	@echo "Checking file headers"
	@python3 test/verify_boilerplate.py

