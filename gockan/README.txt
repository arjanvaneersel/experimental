PACKAGE

package gockan
import "github.com/arjanvaneersel/gockan"

Go language CKAN client and RDF proxy.
Written by the Open Knowledgeg Foundation in 2011.
Distributed under the terms of the GPL version 3 or later.

Work in Progress.

Installation
============

As of 2011/03/29 this requires the weekly snapshot of Go.
On ubuntu linux you can get this by doing:

    add-apt-repository ppa:niemeyer/ppa
    apt-get update
    apt-get install golang-weekly

It also requires the raptor library version 2. There may
be packages for this, or else installing it from source
is straightforward. http://librdf.org/raptor/

For further instructions see the documentation in the
ckand subdirectory.


CONSTANTS

const Version = "0.1"


VARIABLES

var UserAgent string


FUNCTIONS

func Catalog(repo Repository, url, pkgbase string) (ch chan *goraptor.Statement)
Produce an rdf representation of the top level of a dcat:Catalog
from the given repository using the given package base.

func CatalogRecord(pkg *model.Package, url string) (ch chan *goraptor.Statement)
Produce an RDF representation of the given package as a
dcat:CatalogRecord using the given url as package identifier


TYPES

type JsonClient struct {
    ApiBase, ApiKey, ExtraBase, TagBase, GroupBase string
    // contains unexported fields
}

func (jc *JsonClient) Count() int

func (jc *JsonClient) GetPackage(id string) (pkg *model.Package, err error)

func (jc *JsonClient) Packages() (ch chan *model.Package, err error)

func (jc *JsonClient) PutPackage(pkg *model.Package) (err error)

type MemRepo struct {
    // contains unexported fields
}

func (mr *MemRepo) Count() int

func (mr *MemRepo) GetPackage(id string) (pkg *model.Package, err error)

func (mr *MemRepo) Packages() (ch chan *model.Package, err error)

func (mr *MemRepo) PutPackage(pkg *model.Package) (err error)

type Repository interface {
    // number of packages present
    Count() int
    // a channel over which all packages are sent
    Packages() (chan *model.Package, error)
    // retrieve a particular package
    GetPackage(id string) (pkg *model.Package, err error)
    // replace a particular package
    PutPackage(pkg *model.Package) (err error)
}
Interface to the Repository

func NewJsonClient(args ...string) (repo Repository)

func NewMemRepo() (repo Repository)


SUBDIRECTORIES

	.hg
	ckand
	mime
	model
	rdf
	store
