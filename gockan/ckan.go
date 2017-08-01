package gockan

import "github.com/arjanvaneersel/gockan/model"

const Version = "0.1.1"
const UserAgent = "go CKAN v0.1.1"

// Interface to the Repository
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
