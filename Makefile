GOPATH = $(PWD)

all:
	cd bin && CGO_CFLAGS="-I$(GOPATH)/deps/include" CGO_LDFLAGS="-L$(GOPATH)/deps/libs" go build -gcflags "-N -l" zebra

release:
	cd bin && CGO_CFLAGS="-I$(GOPATH)/deps/include" CGO_LDFLAGS="-L$(GOPATH)/deps/libs" go build -ldflags "-w -s" zebra

tools:
	cd bin && CGO_CFLAGS="-I$(GOPATH)/deps/include" CGO_LDFLAGS="-L$(GOPATH)/deps/libs" go build -ldflags "-w -s" tools

