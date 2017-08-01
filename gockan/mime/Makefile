include $(GOROOT)/src/Make.inc

TARG=github.com/arjanvaneersel/gockan/mime
GOFILES=mime.go

include $(GOROOT)/src/Make.pkg

format:
	gofmt -w *.go

docs:
	gomake clean
	godoc ${TARG} > README.txt
