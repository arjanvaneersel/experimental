include $(GOROOT)/src/Make.inc

TARG=github.com/arjanvaneersel/gockan/model
GOFILES=package.go opmv.go licenses.go

include $(GOROOT)/src/Make.pkg

format:
	gofmt -w *.go

docs:
	gomake clean
	godoc ${TARG} > README.txt
