all:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -gcflags "-N -l" zebra

release:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" zebra

redisprotocol:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" tools/redisprotocol

restore:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" tools/restore

