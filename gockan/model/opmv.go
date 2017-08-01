package model

import (
	"time"

	"github.com/arjanvaneersel/gockan/rdf"
	"bitbucket.org/ww/goraptor"
)

// Implementation of some aspects of the Open Provenance Model
// This structure can be used either persistently or on the fly
type Process struct {
	Node   goraptor.Term
	Used   []goraptor.Term
	OpName goraptor.Term
	Time   time.Time
}

func NewProcess() (proc *Process) {
	proc = &Process{}
	proc.Used = make([]goraptor.Term, 0)
	proc.Time = time.Now().UTC()
	return
}

// Indicate that the process was controlled by the entity
// with the given name
func (proc *Process) SetOp(opname string) {
	name := goraptor.Literal{Value: opname}
	proc.OpName = &name
}

// Indicate that the process used a particular resource
func (proc *Process) Use(url string) {
	uri := goraptor.Uri(url)
	proc.Used = append(proc.Used, &uri)
}

func (proc *Process) ToRdf() (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	switch proc.Node.(type) {
	case *goraptor.Blank:
		proc.Node = rdf.Blank()
	case nil:
		proc.Node = rdf.Blank()
	}
	go func() {
		ch <- &goraptor.Statement{proc.Node, rdf.RDF.U("type"), rdf.OPMV.U("Process"), nil}
		defer close(ch)

		if proc.Used != nil {
			for _, term := range proc.Used {
				ch <- &goraptor.Statement{proc.Node, rdf.OPMV.U("used"), term, nil}
			}
		}

		xsdDateTime := &goraptor.Literal{}
		xsdDateTime.Value = proc.Time.Format(time.RFC3339)
		switch dt := rdf.XSD.U("dateTime").(type) {
		case *goraptor.Uri:
			xsdDateTime.Datatype = string(*dt)
		}
		timestamp := rdf.Blank()
		ch <- &goraptor.Statement{timestamp, rdf.RDF.U("type"), rdf.TIME.U("Instant"), nil}
		ch <- &goraptor.Statement{timestamp, rdf.TIME.U("inXSDDateTime"), xsdDateTime, nil}
		ch <- &goraptor.Statement{proc.Node, rdf.OPMV.U("wasPerformedAt"), timestamp, nil}
		if proc.OpName != nil {
			oper := rdf.Blank()
			ch <- &goraptor.Statement{oper, rdf.RDF.U("type"), rdf.FOAF.U("Agent"), nil}
			ch <- &goraptor.Statement{oper, rdf.FOAF.U("name"), proc.OpName, nil}
			ch <- &goraptor.Statement{proc.Node, rdf.OPMV.U("wasControlledBy"), oper, nil}
		}
	}()
	return ch
}
