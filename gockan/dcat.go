package gockan

import (
	"github.com/arjanvaneersel/gockan/model"
	"github.com/arjanvaneersel/gockan/rdf"
	"bitbucket.org/ww/goraptor"
)

// Produce an rdf representation of the top level of a dcat:Catalog
// from the given repository using the given package base.
func Catalog(repo Repository, url, pkgbase string) (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	cat := goraptor.Uri(url)
	go func() {
		ch <- &goraptor.Statement{&cat, rdf.RDF.U("type"), rdf.DCAT.U("Catalog"), &cat}
		defer close(ch)
		pkgs, err := repo.Packages()
		if err != nil {
			return
		}
		for {
			pkg, ok := <-pkgs
			if !ok {
				break
			}
			resource := goraptor.Uri(pkgbase + pkg.Id)
			ch <- &goraptor.Statement{&cat, rdf.DCAT.U("record"), &resource, &cat}
		}
		alts := rdf.Alternatives(cat)
		for {
			alt, ok := <-alts
			if !ok {
				break
			}
			alt.Graph = &cat
			ch <- alt
		}
	}()
	return
}

// Produce an RDF representation of the given package as a
// dcat:CatalogRecord using the given url as package identifier
func CatalogRecord(pkg *model.Package, url string) (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	rec := goraptor.Uri(url)
	go func() {
		ch <- &goraptor.Statement{&rec, rdf.RDF.U("type"), rdf.DCAT.U("CatalogRecord"), &rec}
		if pkg.Provenance != nil {
			ps := pkg.Provenance.ToRdf()
			for {
				p, ok := <-ps
				if !ok {
					break
				}
				p.Graph = &rec
				ch <- p
			}
			ch <- &goraptor.Statement{&rec, rdf.OPMV.U("wasGeneratedBy"), pkg.Provenance.Node, &rec}
		}
		ps := pkg.ToRdf()
		for {
			p, ok := <-ps
			if !ok {
				break
			}
			p.Graph = &rec
			ch <- p
		}
		ch <- &goraptor.Statement{&rec, rdf.DCAT.U("dataset"), pkg.Node, &rec}

		alts := rdf.Alternatives(rec)
		for {
			alt, ok := <-alts
			if !ok {
				break
			}
			alt.Graph = &rec
			ch <- alt
		}

		close(ch)
	}()

	return
}
