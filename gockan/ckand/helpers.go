package main

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/ww/goautoneg"
	"bitbucket.org/ww/goraptor"
)

// helper function for autonegotiation. this looks at the requested
// file extension if it exists and uses that type so long as it is
// compatible with the accept header. otherwise accept header
// negotiation takes place as normal
func negotiate(req *http.Request) (resource, content_type, format string) {
	resource = req.URL.Path

	accept := req.Header.Get("Accept")
	n := strings.LastIndex(resource, ".")
	switch {
	case len(accept) == 0: // assume json
		format = "json"
		content_type = "application/json"
	case n > 0 && n < len(resource)-1:
		format = resource[n+1:]
		resource = resource[:n]
		if format == "json" {
			alternatives := []string{"application/json", "text/javascript"}
			content_type = goautoneg.Negotiate(accept, alternatives)
		} else {
			syntax := goraptor.SerializerSyntax[format]
			if syntax != nil {
				// negotiate here to make sure what was requested is actually
				// acceptable according to the provided header
				content_type = goautoneg.Negotiate(accept, []string{syntax.MimeType})
				// and because there are different syntax names for the same
				// mime type, e.g. rdfxml, rdfxml-abbrev
			}
		}
	default:
		content_type = goautoneg.Negotiate(accept, mime_types)
		format = mime_map[content_type]
	}
	return
}

func servererror(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	body := fmt.Sprintf(`
500 Internal Server Error

Oops something went badly wrong:

%s
`, err)
	w.Write([]byte(body))
	return
}

func notfound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	body := `
404 Not Found

Sorry, we were unable to find what you are looking for
`
	w.Write([]byte(body))
	return
}

func unacceptable(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	body := `
406 Not Acceptable Here

The resource you requested was not available in the required format
`
	w.Write([]byte(body))
}
