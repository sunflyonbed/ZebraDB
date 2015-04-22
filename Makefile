GOPATH = $(PWD)
CGO_CFLAGS = -I$(GOPATH)/deps/include
CGO_LDFLAGS = -L$(GOPATH)/deps/libs

all:
	cd bin && go build -gcflags "-N -l" zebra

release:
	cd bin && go build -ldflags "-w -s" zebra

