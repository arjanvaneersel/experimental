PACKAGE

package rdf
import "github.com/arjanvaneersel/gockan/rdf"

RDF helper routines.


FUNCTIONS

func Alternatives(uri goraptor.Uri) (ch chan *goraptor.Statement)
For the given URI, generate triples of the form,

   <uri> foaf:isPrimaryTopicOf <uri.ext> .
   <uri.ext> a foaf:Document;
       foaf:primaryTopic <uri>;
       dc:format [ a dc:IMT; rdfs:label "ext"; rdf:value "mime_type" ].

Where ext and mime_type come from those that are supported
by the raptor serializer

func Blank() goraptor.Term
Generate a new blank node

func Serializer(format string) (serializer *goraptor.Serializer)
Return a serializer, basically a wrapper around the
goraptor.NewSerializer that also takes care of setting the
namespaces that we know about.


TYPES

type Namespace string
An RDF namespace

var DC Namespace
Dublin Core

var DCAT Namespace
Data Catalog

var FOAF Namespace
Friend of a Friend

var LIC Namespace
Licenses

var MOAT Namespace
Meaning of a Tag

var OPMV Namespace
Open Provenance Model Vocabulary

var RDF Namespace
RDF

var RDFS Namespace
RDF Schema

var REV Namespace
Reviews

var TIME Namespace
OWL Time

var XSD Namespace
XML Schema Datatypes

func (ns Namespace) U(value string) goraptor.Term
Return a goraptor.Term (actually &goraptor.Uri) that
is the concatenation of the namespace and the provided
value.

