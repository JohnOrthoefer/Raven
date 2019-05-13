TARGET=nmapparse
SRC=$(TARGET).go version.go
OBJ=$(TARGET).so
BIN=$(TARGET)
INI=base.ini

$(BINDIR)/$(BIN): $(BIN) $(INI)
	cp $(BIN) $(BINDIR)
	cp $(INI) $(BINDIR)

version.go:
	$(GENVER)

$(BIN): $(SRC)
	$(GO) build $(SRC)

clean:
	rm -f $(BIN) .last_commit version.go
