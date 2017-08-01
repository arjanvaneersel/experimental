PACKAGE

package model
import "github.com/arjanvaneersel/gockan/model"

VARIABLES

var Licenses map[string]string
A mapping of license identifiers to descriptions


TYPES

type Agent struct {
    // URI or blank node identifying this agent
    Node goraptor.Term
    // name as in foaf:name
    Name string
    // mailbox as in foaf:mbox but without the leading mailto:
    Mbox string
    // homepage
    Homepage string
}
An agent that typically has some relationship to a
package. This can be a maintainer or an author and
is the homologue of a foaf:Agent

func (agent *Agent) ToRdf() (ch chan *goraptor.Statement)

type Package struct {
    Node                          goraptor.Term
    ExtraBase, TagBase, GroupBase string
    Provenance                    *Process
    Id                            string
    Revision                      string
    State                         string
    Created                       *time.Time
    Modified                      *time.Time
    Homepage                      goraptor.Term
    Name                          string
    Title                         string
    Version                       string
    Author                        *Agent
    Maintainer                    *Agent
    License                       string
    Tags                          []string
    Groups                        []string
    Extras                        map[string]string
    Resources                     []*Resource
    RatingsCount                  int
    RatingsAvg                    float64
    Notes, NotesRendered          string
    Relationships                 []interface{}
}
Package structure that mirrors that present in CKAN

func PackageFromMap(pkgmap map[string]interface{}) (pkg *Package)
Create a package from a map. This will typically be called with the
results of json.Unmarshal

func (pkg *Package) ToMap() (pkgmap map[string]interface{})

func (pkg *Package) ToRdf() (ch chan *goraptor.Statement)

type Process struct {
    Node   goraptor.Term
    Used   []goraptor.Term
    OpName goraptor.Term
    Time   *time.Time
}
Implementation of some aspects of the Open Provenance Model
This structure can be used either persistently or on the fly

func NewProcess() (proc *Process)

func (proc *Process) SetOp(opname string)
Indicate that the process was controlled by the entity
with the given name

func (proc *Process) ToRdf() (ch chan *goraptor.Statement)

func (proc *Process) Use(url string)
Indicate that the process used a particular resource

type Resource struct {
    Node        goraptor.Term
    Id          string
    Package     string
    Group       string
    AccessURL   string
    Format      string
    Hash        string
    Description string
    Position    int
}
Resource mirroring package resources in CKAN

func ResourceFromMap(resmap map[string]interface{}) (res *Resource)
Create a resource from a map representation. This will
typically be called with the results of json.Unmarshal

func (res *Resource) ToMap() (resmap map[string]interface{})
Return a map representation of the resource suitable for
passing to json.Marshal

func (res *Resource) ToRdf() (ch chan *goraptor.Statement)

