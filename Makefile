GO := go
PKGS = $(shell $(GO) list ./... | grep -v /vendor/)

ARCH ?= $(shell go env GOARCH)
OS ?= $(shell go env GOOS)

OUT_DIR := build
OUTBIN_DIR := ${OUT_DIR}/bin
SRC_PREFIX := github.com/ponypaver/docker-tunnel/cmd
PKG_PREFIX := github.com/ponypaver/docker-tunnel/pkg

ALL_TARGETS := dockertunnel
ALL_TARGETS_WITH_OS := $(addsuffix -os-%,$(ALL_TARGETS))

# compile with default platform options
# eg: make
$(ALL_TARGETS):
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 $(GO) build -ldflags "${LDFLAGS}" -o $(OUTBIN_DIR)/$(OS)/$@ $(SRC_PREFIX)/$@

# cross compile
# eg: make dockertunnel-os-linux
$(ALL_TARGETS_WITH_OS):
	@$(MAKE) OS=$* $(firstword $(subst -os-, ,$@))

build: all

test:
	@$(GO) test $(PKGS)
.PHONY: test

clean:
	@rm -rf ${OUT_DIR}/*
.PHONY: clean
