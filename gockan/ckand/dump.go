package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arjanvaneersel/gockan"
	"github.com/arjanvaneersel/gockan/model"
	"github.com/arjanvaneersel/gockan/rdf"
	"bitbucket.org/ww/goraptor"
	"github.com/vuleetu/goconfig/config"
)

// Dump the contents of the working set out to disk in nquads format.
// Puts dumps in directories named with the date and manages a link
// to the most recent.
func Dump(cfg *config.Config) (err error) {
	cat_uri, err := cfg.String("urls", "catalogue")
	if err != nil {
		return
	}
	pkg_base, err := cfg.String("urls", "package_base")
	if err != nil {
		return
	}
	dump_dir, err := cfg.String("urls", "dump_directory")
	if err != nil {
		dump_dir = "dumps"
		err = nil
	}

	serializer := rdf.Serializer("nquads")
	defer serializer.Free()

	now := time.Now().UTC()
	dump_subdir := fmt.Sprintf("%s/%04d/%02d", dump_dir, now.Year, now.Month)
	err = os.MkdirAll(dump_subdir, 0755)
	filename := fmt.Sprintf("%s/catalogue-%s.nquads", dump_subdir, now.Format(time.RFC3339))

	fp, err := os.Open(filename) //, os.O_WRONLY|os.O_CREAT|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer fp.Close()

	log.Printf("dumping to %s", filename)

	serializer.SetFile(fp, "")

	ch := make(chan *goraptor.Statement)
	go func() {
		defer close(ch)
		catalogue := goraptor.Uri(cat_uri)

		log.Print("generating catalogue")
		cat := gockan.Catalog(repo, string(catalogue), pkg_base)
		for {
			s, ok := <-cat
			if !ok {
				break
			}
			s.Graph = &catalogue
			ch <- s
		}

		log.Print("processing packages")
		packages, err := repo.Packages()
		if err != nil {
			return
		}
		for {
			pkg, ok := <-packages
			if !ok {
				break
			}
			pkguri := goraptor.Uri(pkg_base + pkg.Id)
			pkgrdf := gockan.CatalogRecord(pkg, string(pkguri))
			for {
				s, ok := <-pkgrdf
				if !ok {
					break
				}
				s.Graph = &pkguri
				ch <- s
			}
		}

		log.Print("adding provenance block")
		proc := model.NewProcess()
		proc.SetOp(server_software)
		statements := proc.ToRdf()
		for {
			s, ok := <-statements
			if !ok {
				break
			}
			s.Graph = &catalogue
			ch <- s
		}
		ch <- &goraptor.Statement{&catalogue, rdf.OPMV.U("wasGeneratedBy"), proc.Node, &catalogue}
	}()

	serializer.AddN(ch)

	current := dump_dir + "/catalogue-current.nquads"
	_ = os.Remove(current)

	err = os.Link(filename, current)

	return
}
