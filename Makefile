BINDIR=bin

#.PHONY: pbs

all: m a i test
#
#pbs:
#	cd pbs/ && $(MAKE)
#

tp:=./

test:
	go build  -ldflags '-w -s' -o $(BINDIR)/ctest mac/*.go
m:
	CGO_CFLAGS=-mmacosx-version-min=10.11 \
	CGO_LDFLAGS=-mmacosx-version-min=10.11 \
	GOARCH=amd64 GOOS=darwin go build  --buildmode=c-archive -o $(BINDIR)/simple.a mac/*.go
	cp mac/callback.h $(BINDIR)/
a:
	 gomobile bind -v -o $(BINDIR)/simple.aar -target=android -ldflags=-s github.com/redeslab/simple-go/android
i:
	go env -w GOFLAGS=-mod=mod
	gomobile bind -v -o $(BINDIR)/simple.xcframework -target=ios  -ldflags="-w" -ldflags=-s github.com/redeslab/simple-go/ios

clean:
	gomobile clean
	rm $(BINDIR)/*
