include $(GOROOT)/src/Make.inc

TARG=ckand
GOFILES=server.go daemon.go persist.go package.go helpers.go dump.go log.go harvest.go

include $(GOROOT)/src/Make.cmd

format:
	gofmt -w *.go

docs:
	gomake clean
	godoc ${TARG} > README.txt
