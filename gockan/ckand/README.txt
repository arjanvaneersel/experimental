PACKAGE

package main
import "github.com/arjanvaneersel/gockan/ckand"

go CKAN Server.

This is an aggregator and protocol translator for use in networks of
CKAN servers. It will eventually provide complete support for calls to
the JSON API (version 2) and provide RDF I/O. The RDF I/O is expressed
primarily in terms of the Data Catalog vocabulary.

INSTALLATION

Installation requires the weekly go release as of 2011-03-15. It also
requires raptor2 from http://librdf.org/raptor2 to be installed
including development headers. With these two prerequisites,
installation of this package procedes:

    goinstall github.com/arjanvaneersel/gockan/ckand

This will clone the mercurial repository and you will end up with the
source code in ${GOROOT}/src/pkg/github.com/arjanvaneersel/gockan and
you can inspect the documentation with, for example,

    godoc github.com/arjanvaneersel/gockan/ckand

The goinstall command will have installed all of the prerequisite packages
but it will not have built the daemon. -- TODO find out why
To build the daemon,

    cd ${GOROOT}/src/pkg/github.com/arjanvaneersel/gockan/ckand
    gomake clean all install

CONFIGURATION

The server has a configuration file in .ini format. An example can be
seen at https://github.com/arjanvaneersel/gockan/src/tip/server/ckand.cfg.
This example configuration can also be found in the source distribution.
The example configuration file contains comments that explain what the
configuration variables do but some notes about the general structure
might be useful here.

The key section operationally is [aggregator]. In this section one or
more sources may be configured. The sources will be combined into one
aggregated instance that is addressed by the package identifiers.

The sources have a type field. Right now the only supported type is
json which means to use the regular CKAN API. There are some
parameters for connecting to the service as well as how URIs are to be
generated for things which have no namespace natively in CKAN's JSON
representation such as tags, groups and particularly extras.

The server takes one command line argument. If it is given -d as an
argument in addition to its config file it will daemonise. It is a
good idea to configure logging to a file in the config when running as
a daemon.

HTTP INTERFACE

The server provides an HTTP interface that pays close attention to
content-type autonegotiation because it will serve different
representations of the same data depending on what was requested. This
behaviour is consistent for GET requests across all of the URLs
served.

To illustrate, consider the resource http://example.org/catalogue that is
served with this software. This resource is available in several formats,
html, rdf/xml, json, turtle amongst others. If one were to request this
resource with an Accept header that looked like this:

     Accept: application/turtle

the result would be a 303 redirect to
http://example.org/catalogue.turtle.  Now the request to this location
will be checked to make sure turtle is really what was asked for and
is consistent with the Accept header, and then it will be returned.

The implications of this are that there are two ways to get a specific
representation of a resource. The first is to set the Accept header
and to follow the redirects, and the second is to append a well known
file extension and put a more permissive Accept header, say,

     Accept: *\/*

as is contained in what is sent by default by most web browsers.

ARCHITECTURE

Internally, the server uses two (or more) repository types. The
primary storage is gockan.MemRepo, an in-memory map of package
identifiers to packages. The remote CKAN instances are represented by
gockan.JsonRepo and the remote repositories are periodically polled
for updates to feed into the main memory repository. It may sound
dodgy to try to keep everything in memory, however the implementation
is quite memory efficient and there is not, at present, a substantial
amount of metadata being stored in CKAN instances. The memory
footprint of the server including all of this data is only a few dozen
megabytes when aggregating several CKAN instances. In any event, with
a small sacrifice of speed, it is very simple to make a disk based
storage implementation of gockan.Repository using some database.

Support for POST, PUT and DELETE requests, which will simply get
passed on the appropriate remote CKAN instance are not yet supported.
Supporting these operations is not difficult and means that client
software could potentially use a single entry point for updating
records across multiple CKAN instances.

In the case of creating a new resource, another simple to implement
change would be to use a local respository type, such as the in-memory
repository instead of only the remote JSON and directly support the
creation of packages at the aggregation point. This is necessarily
different from updating them because for update operations we know the
source of the records, where to do the update. for a new record we do
not know.


FUNCTIONS

func Dump(cfg *config.Config) (err error)
Dump the contents of the working set out to disk in nquads format.
Puts dumps in directories named with the date and manages a link
to the most recent.

func LoadRepo(repo gockan.Repository, filename string) (modified *time.Time, err error)
Load the contents of the persistent storage into the repository.

XXX persistence in this way needs to be more fully thought out

func LogRequest(req *http.Request)
log the supplied http request

func MirrorRepo(src, dst gockan.Repository, filename string) (err error)
Mirror the source repository into the destination repository

XXX persistence in this way needs to be more fully thought out

func PackageList(w http.ResponseWriter, req *http.Request)
Handler produces a list of packages for JSON or dcat:Catalogue for RDF

func PackageService(w http.ResponseWriter, req *http.Request)
Handler produces a JSON rendering of a package or a dcat:CatalogRecord
for RDF requests

func SignalHandler(cfg *config.Config)
Handle certain signals gracefully.

    SIGINT and SIGTERM terminate the program
    SIGHUP causes the program to restart
    SIGINFO causes the program to print status information in the log
    SIGUSR1 causes the program to dump the contents of its working set


TYPES

type Harvester struct {
    // contains unexported fields
}

func NewHarvester(cfg *config.Config, name string, src, dst gockan.Repository) *Harvester

func (h *Harvester) Start()

