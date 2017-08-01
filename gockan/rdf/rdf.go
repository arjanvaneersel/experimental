/*
RDF helper routines.
*/
package rdf

import (
	"bitbucket.org/ww/goraptor"
	"fmt"
	"strings"
	"sync"
)

// An RDF namespace
type Namespace string

// Return a goraptor.Term (actually &goraptor.Uri) that
// is the concatenation of the namespace and the provided
// value.
func (ns Namespace) U(value string) goraptor.Term {
	uri := goraptor.Uri(string(ns) + value)
	return &uri
}

// RDF
var RDF Namespace
// RDF Schema
var RDFS Namespace
// Dublin Core
var DC Namespace
// Data Catalog
var DCAT Namespace
// Friend of a Friend
var FOAF Namespace
// Licenses
var LIC Namespace
// Meaning of a Tag
var MOAT Namespace
// Open Provenance Model Vocabulary
var OPMV Namespace
// OWL Time
var TIME Namespace
// Reviews
var REV Namespace
// XML Schema Datatypes
var XSD Namespace

var blank_seq uint64
var blank_mutex sync.Mutex

// Generate a new blank node
func Blank() goraptor.Term {
	blank_mutex.Lock()
	bnode := goraptor.Blank(fmt.Sprintf("genid%d", blank_seq))
	blank_seq += 1
	blank_mutex.Unlock()
	return &bnode
}

// For the given URI, generate triples of the form,
//
//    <uri> foaf:isPrimaryTopicOf <uri.ext> .
//    <uri.ext> a foaf:Document;
//        foaf:primaryTopic <uri>;
//        dc:format [ a dc:IMT; rdfs:label "ext"; rdf:value "mime_type" ].
//
// Where ext and mime_type come from those that are supported
// by the raptor serializer  
func Alternatives(uri goraptor.Uri) (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	go func() {
		for format, syntax := range goraptor.SerializerSyntax {
			/// XXX this is kludgy... why do we strip these out?
			if strings.Index(format, ".") != -1 {
				continue
			}
			doc := goraptor.Uri(string(uri) + "." + format)
			ch <- &goraptor.Statement{&doc, RDF.U("type"), FOAF.U("Document"), nil}
			ch <- &goraptor.Statement{&doc, FOAF.U("primaryTopic"), &uri, nil}
			dcformat := Blank()
			ch <- &goraptor.Statement{dcformat, RDF.U("type"), DC.U("IMT"), nil}
			label := goraptor.Literal{Value: syntax.Name}
			ch <- &goraptor.Statement{dcformat, RDFS.U("label"), &label, nil}
			value := goraptor.Literal{Value: syntax.MimeType}
			ch <- &goraptor.Statement{dcformat, RDF.U("value"), &value, nil}
			ch <- &goraptor.Statement{&doc, DC.U("format"), dcformat, nil}
			ch <- &goraptor.Statement{&uri, FOAF.U("isPrimaryTopicOf"), &doc, nil}
		}
		close(ch)
	}()
	return
}

// Return a serializer, basically a wrapper around the 
// goraptor.NewSerializer that also takes care of setting the
// namespaces that we know about.
func Serializer(format string) (serializer *goraptor.Serializer) {
	serializer = goraptor.NewSerializer(format)
	serializer.SetNamespace("rdf", string(RDF))
	serializer.SetNamespace("rdfs", string(RDFS))
	serializer.SetNamespace("dc", string(DC))
	serializer.SetNamespace("dcat", string(DCAT))
	serializer.SetNamespace("foaf", string(FOAF))
	serializer.SetNamespace("lic", string(LIC))
	serializer.SetNamespace("moat", string(MOAT))
	serializer.SetNamespace("opmv", string(OPMV))
	serializer.SetNamespace("rev", string(REV))
	serializer.SetNamespace("time", string(TIME))
	serializer.SetNamespace("xsd", string(XSD))
	return
}

func init() {
	RDF = Namespace("http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	RDFS = Namespace("http://www.w3.org/2000/01/rdf-schema#")
	DC = Namespace("http://purl.org/dc/terms/")
	DCAT = Namespace("http://www.w3.org/ns/dcat#")
	FOAF = Namespace("http://xmlns.com/foaf/0.1/")
	LIC = Namespace("http://purl.org/okfn/licenses/")
	MOAT = Namespace("http://moat-project.org/ns#")
	OPMV = Namespace("http://purl.org/net/opmv/ns#")
	TIME = Namespace("http://www.w3.org/2006/time#")
	REV = Namespace("http://purl.org/stuff/rev#")
	XSD = Namespace("http://www.w3.org/2001/XMLSchema#")
}
