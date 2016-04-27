export CGO_ENABLED=0

BUILDDIR = ./.build

.PHONY: all build 

build:
	mkdir -p $(BUILDDIR)
	go build -i -o $(BUILDDIR)/cloudthing .

clean:
	rm -Rf $(BUILDDIR)

all: build
	

