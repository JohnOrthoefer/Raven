ifndef ROOTDIR
$(error Please run from package root)
endif

TARGET=huginn
VERSION=$(shell pwd)/license/version.go
SRC=$(TARGET).go
OBJ=$(TARGET).so
BIN=$(TARGET)

all: $(BINDIR)/$(BIN)
	$(MAKE) -C raven
	rsync -av templates/ $(TMPLDIR)

$(BINDIR)/$(BIN): $(BIN) 
	cp $(BIN) $(BINDIR)

$(VERSION):
	$(GENVER) $(shell dirname $(VERSION))

$(BIN): $(SRC) $(VERSION)
	$(GO) build $(SRC)

clean:
	rm -f $(BIN) .last_commit $(VERSION)
	$(MAKE) -C raven clean
