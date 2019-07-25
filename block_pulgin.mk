BASEDIR := $(shell pwd)
SOURCEDIR= contrib/drivers/volume/plugins
BUILDDIR= build/out/lib/block
BUILDTEMPDIR=build/temp/lib/block

SOURCES = $(shell find $(SOURCEDIR) -name '*.go')

OBJECTS = $(patsubst $(SOURCEDIR)/%.go, $(BUILDTEMPDIR)/%.so, $(SOURCES))

all: prehandle buildso
	@find $(BUILDTEMPDIR) -mindepth 2 -type f -name *.so| xargs -i{} cp {} $(BUILDDIR)

prehandle:
	@mkdir $(BUILDDIR) -p

buildso:$(OBJECTS)
$(BUILDTEMPDIR)/%.so: $(SOURCEDIR)/%.go
	go build -ldflags '-w -s' -buildmode=plugin -o $@ $<
clean: 
	rm $(BUILDTEMPDIR) $(BUILDDIR) -rf
