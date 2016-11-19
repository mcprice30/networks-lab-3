# Constants for running go.
SETPATH=GOPATH="$(shell pwd)"
GOCMD=bin/go/bin/go
GO=$(SETPATH) $(GOCMD)

# Go setup flags
GO_DOWNLOAD_SRC="https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz"
GO_TARBALL=godownload.tar.gz
GROOT=$(shell pwd)/bin/go
export GOROOT:=$(GROOT)

all:  clean slave master
	@echo "Done!"

setup: setup-dirs
	@echo "Downloading Go..."
	@curl $(GO_DOWNLOAD_SRC) > bin/$(GO_TARBALL) && \
	echo "Extracting Go..." && \
	tar -xzf bin/$(GO_TARBALL) -C bin && \
	echo "Go installed successfully!"

setup-dirs:
	@mkdir -p bin

fmt:
	@$(GO) fmt loader memory parser simulator util


slave:
	cc src/Slave.c -o slave -lpthread

master:
	$(GO) build src/master.go

udisplay:
	cc provided/UDPServerDisplay.c -o udisplay
tdisplay:
	cc provided/TCPServerDisplay.c -o tdisplay

clean:
	rm -f slave master udisplay tdisplay
