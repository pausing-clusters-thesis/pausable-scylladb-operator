all: build
.PHONY: all

SHELL :=/bin/bash -euEo pipefail -O inherit_errexit

IMAGE_TAG ?= latest
IMAGE_REF ?= docker.io/rzetelskik/pausable-scylladb-operator:$(IMAGE_TAG)
IMAGE_CONTAINERFILE ?=./images/pausable-scylladb-operator/Containerfile

GO_PACKAGE ?=github.com/pausing-clusters-thesis/pausable-scylladb-operator
GO_MODULE ?=github.com/pausing-clusters-thesis/pausable-scylladb-operator

MAKE_REQUIRED_MIN_VERSION:=4.2 # for SHELLSTATUS

# Support container build from git worktrees where the parent git folder isn't available.
GIT ?=git

GIT_TAG ?=$(shell [ ! -d ".git/" ] || $(GIT) describe --long --tags --abbrev=7 --match 'v[0-9]*')$(if $(filter $(.SHELLSTATUS),0),,$(error $(GIT) describe failed))
GIT_TAG_SHORT ?=$(shell [ ! -d ".git/" ] || $(GIT) describe --tags --abbrev=7 --match 'v[0-9]*')$(if $(filter $(.SHELLSTATUS),0),,$(error $(GIT) describe failed))
GIT_COMMIT ?=$(shell [ ! -d ".git/" ] || $(GIT) rev-parse --short "HEAD^{commit}" 2>/dev/null)$(if $(filter $(.SHELLSTATUS),0),,$(error $(GIT) rev-parse failed))
GIT_TREE_STATE ?=$(shell ( ( [ ! -d ".git/" ] || $(GIT) diff --quiet ) && echo 'clean' ) || echo 'dirty')

GO ?=go
GO_MODULE ?=$(shell $(GO) list -m)$(if $(filter $(.SHELLSTATUS),0),,$(error failed to list go module name))
GOPATH ?=$(shell $(GO) env GOPATH)
GOOS ?=$(shell $(GO) env GOOS)
GOEXE ?=$(shell $(GO) env GOEXE)
GOFMT ?=gofmt
GOFMT_FLAGS ?=-s -l

GO_VERSION :=$(shell $(GO) version | sed -E -e 's/.*go([0-9]+.[0-9]+.[0-9]+).*/\1/')
GO_PACKAGE ?=$(shell $(GO) list -m -f '{{ .Path }}' || echo 'no_package_detected')
GO_PACKAGES ?=./...

go_packages_dirs :=$(shell $(GO) list -f '{{ .Dir }}' $(GO_PACKAGES) || echo 'no_package_dir_detected')
GO_TEST_PACKAGES ?=$(GO_PACKAGES)
GO_BUILD_PACKAGES ?=./cmd/...
GO_BUILD_PACKAGES_EXPANDED ?=$(shell $(GO) list $(GO_BUILD_PACKAGES))
go_build_binaries =$(notdir $(GO_BUILD_PACKAGES_EXPANDED))
GO_BUILD_FLAGS ?=-trimpath
GO_BUILD_BINDIR ?=
GO_LD_EXTRA_FLAGS ?=
GO_TEST_PACKAGES :=./pkg/... ./cmd/...
GO_TEST_FLAGS ?=-race
GO_TEST_COUNT ?=
GO_TEST_EXTRA_FLAGS ?=
GO_TEST_ARGS ?=
GO_TEST_EXTRA_ARGS ?=
GO_TEST_E2E_EXTRA_ARGS ?=

JQ ?=jq
YQ ?=yq -e

CODEGEN_PKG ?=./vendor/k8s.io/code-generator
CODEGEN_HEADER_FILE ?=/dev/null

api_groups :=$(patsubst %/,%,$(wildcard ./pkg/api/*/))

api_package_dirs :=$(api_groups)
api_packages =$(call expand_go_packages_with_spaces,$(addsuffix /...,$(api_package_dirs)))

CONTROLLER_GEN ?=$(GO) run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen --
CRD_FILES ?=$(shell find ./pkg/api/ -name '*.yaml')$(if $(filter $(.SHELLSTATUS),0),,$(error "can't find CRDs"))

define version-ldflags
-X $(1).versionFromGit="$(GIT_TAG)" \
-X $(1).commitFromGit="$(GIT_COMMIT)" \
-X $(1).gitTreeState="$(GIT_TREE_STATE)" \
-X $(1).buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"
endef
GO_LD_FLAGS ?=-ldflags '$(strip $(call version-ldflags,$(GO_PACKAGE)/pkg/version) $(GO_LD_EXTRA_FLAGS))'

export DOCKER_BUILDKIT :=1
export GOVERSION :=$(shell go version)
export KUBEBUILDER_ASSETS :=$(GOPATH)/bin
export PATH :=$(GOPATH)/bin:$(PATH):

# $1 - required version
# $2 - current version
define is_equal_or_higher_version
$(strip $(filter $(2),$(firstword $(shell printf '%s\n%s' '$(1)' '$(2)' | sort -V -r -b))))
endef

# $1 - program name
# $2 - required version variable name
# $3 - current version string
define require_minimal_version
$(if $($(2)),\
$(if $(strip $(call is_equal_or_higher_version,$($(2)),$(3))),,$(error `$(1)` is required with minimal version "$($(2))", detected version "$(3)". You can override this check by using `make $(2):=`)),\
)
endef

ifneq "$(MAKE_REQUIRED_MIN_VERSION)" ""
$(call require_minimal_version,make,MAKE_REQUIRED_MIN_VERSION,$(MAKE_VERSION))
endif

# $1 - package name
define build-package
	$(if $(GO_BUILD_BINDIR),mkdir -p '$(GO_BUILD_BINDIR)',)
	$(strip CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) $(GO_LD_FLAGS) \
		$(if $(GO_BUILD_BINDIR),-o '$(GO_BUILD_BINDIR)/$(notdir $(1))$(GOEXE)',) \
	$(1))

endef

# We need to build each package separately so go build creates appropriate binaries
build:
	$(if $(strip $(GO_BUILD_PACKAGES_EXPANDED)),,$(error no packages to build: GO_BUILD_PACKAGES_EXPANDED var is empty))
	$(foreach package,$(GO_BUILD_PACKAGES_EXPANDED),$(call build-package,$(package)))
.PHONY: build

clean:
	$(RM) $(go_build_binaries)
.PHONY: clean

update-gofmt:
	find ./ -name '*.go' -not -path './vendor/*' -exec $(GOFMT) -s -w {} '+'
.PHONY: update-gofmt

# $1 - codegen command
# $2 - extra args
define run-codegen
	$(GO) run "$(CODEGEN_PKG)/cmd/$(1)" --go-header-file='$(CODEGEN_HEADER_FILE)' $(2)

endef

# $1 - api packages
define run-deepcopy-gen
	$(call run-codegen,deepcopy-gen,--output-file='zz_generated.deepcopy.go' $(1))

endef

# $1 - group
# $2 - api packages
# $3 - client dir
define run-client-gen
	$(call run-codegen,client-gen,--clientset-name=versioned --input-base='./' --output-pkg='$(GO_PACKAGE)/$(3)/$(1)/clientset' --output-dir='./$(3)/$(1)/clientset/' $(foreach p,$(2),--input='$(p)'))

endef

# $1 - group
# $2 - api packages
# $3 - client dir
define run-lister-gen
	$(call run-codegen,lister-gen,--output-pkg='$(GO_PACKAGE)/$(3)/$(1)/listers' --output-dir='./$(3)/$(1)/listers/' $(2))

endef

# $1 - group
# $2 - api packages
# $3 - client dir
define run-informer-gen
	$(call run-codegen,informer-gen,--output-pkg='$(GO_PACKAGE)/$(3)/$(1)/informers' --output-dir='./$(3)/$(1)/informers/' --versioned-clientset-package '$(GO_PACKAGE)/$(3)/$(1)/clientset/versioned' --listers-package='$(GO_PACKAGE)/$(3)/$(1)/listers' $(2))

endef

# $1 - packages
expand_go_packages_to_json=$(shell $(GO) list -json $(1) | $(JQ) -sr 'map(select(.Name | test("^v[0-9]+.*$$"))) | reduce .[] as $$item ([]; . + [$$item.ImportPath])')$(if $(filter $(.SHELLSTATUS),0),,$(error failed to expand packages to json: $(1)))

# $1 - packages
expand_go_packages_with_commas=$(shell echo '$(call expand_go_packages_to_json,$(1))' | $(JQ) -r '. | join(",")')$(if $(filter $(.SHELLSTATUS),0),,$(error failed to expand packages with commas: $(1)))

# $1 - packages
expand_go_packages_with_spaces=$(shell echo '$(call expand_go_packages_to_json,$(1))' | $(JQ) -r '. | join(" ")')$(if $(filter $(.SHELLSTATUS),0),,$(error failed to expand packages with spaces: $(1)))

# $1 - group
# $2 - api packages
# $3 - client dir
define run-client-generators
	$(call run-client-gen,$(1),$(2),$(3))
	$(call run-lister-gen,$(1),$(2),$(3))
	$(call run-informer-gen,$(1),$(2),$(3))

endef

define run-update-codegen
	$(call run-deepcopy-gen,$(addsuffix /...,$(api_groups)))
	$(foreach group,$(api_groups),$(call run-client-generators,$(notdir $(group)),$(call expand_go_packages_with_spaces,$(group)/...),pkg/client))

endef

update-codegen:
	$(call run-update-codegen)
.PHONY: update-codegen

# $1 - api package
# $2 - output dir
# We need to cleanup `---` in the yaml output manually because it shouldn't be there and it breaks opm.
define run-crd-gen
	$(CONTROLLER_GEN) crd paths='$(1)' output:dir='$(2)'
	find '$(2)' -mindepth 1 -maxdepth 1 -type f -name '*.yaml' -exec $(YQ) -i eval '... style=""' {} \;

endef

# $1 - dir prefix
define generate-crds
	$(foreach p,$(api_packages),$(call run-crd-gen,$(subst $(GO_MODULE)/,,./$(p)),$(1)$(subst $(GO_MODULE)/,,./$(p))))

endef

update-crds:
	$(call generate-crds,)
.PHONY: update-crds

update: update-crds update-codegen update-gofmt
.PHONY: update

help:
	$(info The following make targets are available:)
	@$(MAKE) -f $(firstword $(MAKEFILE_LIST)) --print-data-base --question no-such-target 2>&1 | grep -v 'no-such-target' | \
	grep -v -e '^no-such-target' -e '^makefile' | \
	awk '/^[^.%][-A-Za-z0-9_]*:/	{ print substr($$1, 1, length($$1)-1) }' | sort -u
.PHONY: help

image:
	buildah build --file "$(IMAGE_CONTAINERFILE)" --tag "$(IMAGE_REF)" .
.PHONY: image