PROJECT?=must-override

.PHONY: all
all: build test

.PHONY: gazell-update
gazelle-update:
	bazel run //:gazelle -- update

.PHONY: gazell-update-mod
gazelle-update-mod:
	GO111MODULE=on go mod tidy
	bazel run //:gazelle -- update-repos -from_file=go.mod
	bazel run //:gazelle -- update

.PHONY: build
build:
	bazel build //...

.PHONY: test
test:
	bazel test --test_verbose_timeout_warnings //...

.PHONY: clean
clean:
	bazel clean

.PHONY: gen
gen:
	bazel-bin/darwin_amd64_stripped/gke-tf gen -d /tmp/tf/ -f examples/example.yaml -p ${PROJECT}
