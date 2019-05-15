ifndef ROOTDIR
$(error Please run from package root)
endif

TARGET=huginn
SRC=$(TARGET).go
OBJ=$(TARGET).so
BIN=$(TARGET)

all: $(BINDIR)/$(BIN)
	$(MAKE) -C raven
	rsync -av templates/ $(TMPLDIR)

$(BINDIR)/$(BIN): $(BIN) 
	cp $(BIN) $(BINDIR)

licenses/version.go: 
	$(GENVER) ${PWD}/licenses

$(BIN): $(SRC)
	$(GO) build $(SRC)

clean:
	rm -f $(BIN) .last_commit version.go
	$(MAKE) -C raven clean
