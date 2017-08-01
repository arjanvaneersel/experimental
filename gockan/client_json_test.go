package gockan

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/arjanvaneersel/gockan/rdf"
)

func TestGetPackage(t *testing.T) {
	repo := NewJsonClient()
	pkg, err := repo.GetPackage("438bbe0e-4f2a-4021-b238-a34e5bf31c74")
	if err != nil {
		t.Error(err)
	}

	repo.PutPackage(pkg)

	serializer := rdf.Serializer("turtle")
	defer serializer.Free()

	statements := pkg.ToRdf()
	str, err := serializer.Serialize(statements, "")
	if err != nil {
		t.Error(err)
	}
	log.Print(str)

	buf := bytes.NewBuffer(make([]byte, 0, 16384))
	enc := gob.NewEncoder(buf)
	enc.Encode(pkg)
	data := buf.Bytes()
	log.Print(data)
	log.Print(len(data))
}
