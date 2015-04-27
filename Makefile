all: zebra redisprotocol restore save

zebra:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -gcflags "-N -l" zebra

redisprotocol:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" tools/redisprotocol

restore:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" tools/restore

save:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" tools/save

release:
	cd bin && GOPATH=$(PWD) CGO_CFLAGS="-I$(PWD)/deps/include" CGO_LDFLAGS="-L$(PWD)/deps/libs" go build -ldflags "-w -s" zebra

