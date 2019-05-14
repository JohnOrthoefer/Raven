export VER=$(shell git tag | head -1)
export ROOTDIR=$(PWD)
export INSTALLDIR=${ROOTDIR}/srv
export BINDIR=$(INSTALLDIR)/bin
export PLUGINDIR=$(BINDIR)/plugins
export ETCDIR=$(INSTALLDIR)/etc
export TMPLDIR=$(INSTALLDIR)/templates

GO := $(shell which go)
ifneq ($(.SHELLSTATUS),0) 
	$(error No go compiler in path)
endif
export GO

all: $(BINDIR) $(PLUGINDIR) $(TMPLDIR) $(ETCDIR)
	@echo Version- $(VER)
	$(MAKE) -C src all

$(BINDIR):
	mkdir -p $(BINDIR)
$(ETCDIR):
	mkdir -p $(ETCDIR)
$(PLUGINDIR):
	mkdir -p $(PLUGINDIR)
$(TMPLDIR):
	mkdir -p $(TMPLDIR)

depend:
	$(GO) get github.com/go-ini/ini

clean:
	rm -rf $(INSTALLDIR)
	$(MAKE) -C src clean

tar: clean
	-mkdir ../raven-$(VER)
	rsync -av src Makefile COPYING README.md docs etc ../raven-$(VER)/
	(cd ..; tar czf raven-$(VER).tar.gz raven-$(VER))

dist-deb:
	bzr dh-make raven $(VER) raven-$(VER).tar.gz 
