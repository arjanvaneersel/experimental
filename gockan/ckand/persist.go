package main

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"

	"github.com/arjanvaneersel/gockan"
	"github.com/arjanvaneersel/gockan/model"
)

// Mirror the source repository into the destination repository
//
// XXX persistence in this way needs to be more fully thought out
func MirrorRepo(src, dst gockan.Repository, filename string) (err error) {
	ch, err := src.Packages()
	if err != nil {
		return
	}

	tmpfile := filename + ".tmp"
	_ = os.Remove(tmpfile)

	fp, err := os.Open(tmpfile) //, os.O_WRONLY|syscall.O_CREAT|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer fp.Close()

	enc := gob.NewEncoder(fp)
	err = enc.Encode(time.Now().UTC())

	for {
		pkg, ok := <-ch
		if !ok {
			break
		}

		err = dst.PutPackage(pkg)
		if err != nil {
			log.Print(err, pkg)
		}

		err = enc.Encode(pkg)
		if err != nil {
			log.Print(err, pkg)
			err = nil
		}
	}

	_ = os.Remove(filename)
	err = os.Link(tmpfile, filename)
	if err != nil {
		return
	}
	err = os.Remove(tmpfile)
	return
}

// Load the contents of the persistent storage into the repository.
//
// XXX persistence in this way needs to be more fully thought out
func LoadRepo(repo gockan.Repository, filename string) (modified *time.Time, err error) {
	fp, err := os.Open(filename) //, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer fp.Close()

	dec := gob.NewDecoder(fp)

	err = dec.Decode(modified)
	if err != nil {
		return
	}

	for {
		pkg := &model.Package{}
		err = dec.Decode(pkg)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		err = repo.PutPackage(pkg)
		if err != nil {
			return
		}
	}
	return
}
