include $(GOPATH)/src/github.com/arjanvaneersel/gockan/Make.inc

TARG=github.com/arjanvaneersel/gockan
GOFILES=ckan.go repo.go client_json.go docs.go dcat.go

include $(GOPATH)/src/github.com/arjanvaneersel/gockan/Make.pkg

format:
	gofmt -w *.go

docs:
	gomake clean
	godoc ${TARG} > README.txt
