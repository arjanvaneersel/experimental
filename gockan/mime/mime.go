/*
CKAN format tags to mime types.
*/
package mime

type MimeType struct {
	// The human readable friendly label
	Label string
	// The actual mime type
	Value string
}
var MimeTypes map[string]MimeType

func init() {
	MimeTypes = make(map[string]MimeType)
	MimeTypes["json"] = MimeType{"JSON", "application/json"}
	MimeTypes["text/n3"] = MimeType{"N3", "text/n3"}
	MimeTypes["application/rdf+xml"] = MimeType{"RDF/XML", "application/rdf+xml"}
	MimeTypes["example/n3"] = MimeType{"N3", "text/n3"}
	MimeTypes["example/ntriples"] = MimeType{"NT", "text/plain"}
	MimeTypes["example/rdf+xml"] = MimeType{"RDF/XML", "application/rdf+xml"}
	MimeTypes["api/sparql"] = MimeType{"SPARQL", "application/sparql-results+xml"}
}