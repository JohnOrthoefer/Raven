export VER=$(shell git tag | head -1)
export ROOTDIR=$(PWD)
export GO=/usr/bin/go
export BINDIR=$(ROOTDIR)/bin
export PLUGINDIR=$(BINDIR)/plugins
export TMPLDIR=$(BINDIR)/templates


all: 
	@echo Version- $(VER)
	mkdir -p $(BINDIR)
	mkdir -p $(PLUGINDIR)
	mkdir -p $(TMPLDIR)
	$(MAKE) -C src all

clean:
	rm -rf $(BINDIR)
	$(MAKE) -C src clean

tar:
	-mkdir ../raven-$(VER)
	rsync -av src Makefile COPYING README.md docs etc ../raven-$(VER)/
	(cd ..; tar czf raven-$(VER).tar.gz raven-$(VER))

dist-deb:
	bzr dh-make raven $(VER) raven-$(VER).tar.gz 
