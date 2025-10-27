
VERSION ?= "0.0.1"
GITCOMMIT ?= ""
BuildAt ?= ""
GoOS ?= "linux"
GoArch ?= "amd64"

BIN_DIR = ./bin


all: $(BIN_DIR) build

$(BIN_DIR):
	@echo "Creating directory: $(BIN_DIR)"
	mkdir -p  $(BIN_DIR)

build:
	@echo "build project"
	CGO_ENABLED=0 GOOS=$(GoOS) GOARCH=$(GoArch) go build -ldflags "-X goRelay/pkg.Version=$(VERSION) -X goRelay/pkg.BuildAt=$(BuildAt) -X goRelay/pkg.GitCommit=$(GITCOMMIT)" -o $(BIN_DIR)/goRelay_$(GoOS)_$(GoArch)

clean:
	@echo "cleaning up..."
	rm -f $(BIN_DIR)/goRela*
	rm -f exec_goRelay_*

