ifndef ROOTDIR
$(error Please run from package root)
endif

TARGET=huginn
SRC=$(TARGET).go version.go
OBJ=$(TARGET).so
BIN=$(TARGET)

all: $(BINDIR)/$(BIN)
	$(MAKE) -C raven
	rsync -av templates $(BINDIR)/

$(BINDIR)/$(BIN): $(BIN) 
	cp $(BIN) $(BINDIR)

version.go: 
	$(GENVER)

$(BIN): $(SRC)
	$(GO) build $(SRC)

clean:
	rm -f $(BIN) .last_commit version.go
	$(MAKE) -C raven clean
