ifndef ROOTDIR
$(error Please run from package root)
endif

TARGET=nmapparse
SRC=$(TARGET).go version.go
OBJ=$(TARGET).so
BIN=$(TARGET)
INI=base.ini

ALL: $(BINDIR)/$(BIN) $(ETCDIR)/$(INI)

$(BINDIR)/$(BIN): $(BIN) $(INI)
	cp $(BIN) $(BINDIR)

$(ETCDIR)/$(INI): $(INI)
	cp $(INI) $(ETCDIR)

version.go:
	$(GENVER)

$(BIN): $(SRC)
	$(GO) build $(SRC)

clean:
	rm -f $(BIN) .last_commit version.go
