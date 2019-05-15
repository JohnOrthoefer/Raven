ifndef ROOTDIR
$(error Please run from package root)
endif

TARGET=viassh
SRC=$(TARGET).go
OBJ=$(TARGET).so

$(PLUGINDIR)/$(OBJ): $(OBJ)
	cp $(OBJ) $(PLUGINDIR)

$(OBJ): $(SRC) ../ravenRunCommand/util.go
	$(GO) build -buildmode=plugin $(SRC)
clean:  
	rm -f $(OBJ) 
