TARGET=viassh
SRC=$(TARGET).go
OBJ=$(TARGET).so

$(PLUGINDIR)/$(OBJ): $(OBJ)
	cp $(OBJ) $(PLUGINDIR)

$(OBJ): $(SRC)
	$(GO) build -buildmode=plugin $(SRC)
clean:  
	rm -f $(OBJ) 
