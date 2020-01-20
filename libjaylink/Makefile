# See: https://gitlab.zapb.de/zapb/libjaylink

VER = master
DIR = libjaylink-$(VER)
TGZ = $(DIR).tar.gz  

PATCHFILES = $(sort $(wildcard *.patch ))

PATCH = \
	for f in $(PATCHFILES); do\
	    echo -n $$f": "; \
	    patch -d $(DIR) -p0 -b -z .original < $$f; \
	done

all: untar
	cd $(DIR); ./autogen.sh
	cd $(DIR); ./configure
	$(MAKE) -C $(DIR) $@

.PHONY: untar
untar:
	tar zxf $(TGZ)
	$(PATCH)

.PHONY: clean
clean:
	-rm -rf $(DIR)
