ifndef ROOTDIR
$(error Please run from package root)
endif

ifndef ROOTDIR
$(error Please run from package root)
endif
export GENVER=$(shell pwd)/genver

all: 
	$(MAKE) -C huginn
	$(MAKE) -C nmapparse

clean:
	$(MAKE) -C huginn clean
	$(MAKE) -C nmapparse clean
