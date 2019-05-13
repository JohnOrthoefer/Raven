ifndef ROOTDIR
$(error Please run from package root)
endif

all:
	$(MAKE) -C plugin-nagios
	$(MAKE) -C plugin-ping
	$(MAKE) -C plugin-viassh
clean:
	$(MAKE) -C plugin-nagios clean
	$(MAKE) -C plugin-ping clean
	$(MAKE) -C plugin-viassh clean
