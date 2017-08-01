package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arjanvaneersel/gockan"
	"github.com/arjanvaneersel/gockan/rdf"
)

// factor this out for clarity below, this is the JSON implementation
func package_list() (pkglist []string, err error) {
	pkglist = make([]string, repo.Count())
	ch, err := repo.Packages()
	if err != nil {
		return
	}

	i := 0
	for {
		pkg, ok := <-ch
		if !ok {
			break
		}
		pkglist[i] = pkg.Id
		i++
	}
	return
}

// Handler produces a list of packages for JSON or dcat:Catalogue for RDF
func PackageList(w http.ResponseWriter, req *http.Request) {
	LogRequest(req)
	resource, content_type, format := negotiate(req)
	if len(content_type) == 0 {
		unacceptable(w)
		return
	}

	switch {
	case format == "json":
		pkglist, err := package_list()
		if err != nil {
			servererror(w, err)
			return
		}
		buf, err := json.Marshal(pkglist)
		if err != nil {
			servererror(w, err)
			return
		}
		w.Header().Set("Content-Type", content_type)
		content_length := fmt.Sprintf("%d", len(buf))
		w.Header().Set("Content-Length", content_length)
		if !strings.HasSuffix(req.URL.Path, "."+format) {
			w.Header().Set("Vary", "Accept")
			w.Header().Set("Content-Location", req.URL.Path+".json")
		}
		w.Header().Set("Server", server_software)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	case !strings.HasSuffix(req.URL.Path, "."+format):
		w.Header().Set("Vary", "Accept")
		http.Redirect(w, req, req.URL.Path+"."+format, http.StatusSeeOther)
	default: // rdf
		serializer := rdf.Serializer(format)
		defer serializer.Free()
		resource := "http://" + req.Host + resource
		catalogue := gockan.Catalog(repo, resource, "http://"+req.Host+"/package/")
		data, err := serializer.Serialize(catalogue, "")
		if err != nil {
			servererror(w, err)
			return
		}
		w.Header().Set("Content-Type", content_type)
		content_length := fmt.Sprintf("%d", len(data))
		w.Header().Set("Content-Length", content_length)
		w.Header().Set("Server", server_software)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}

}

// Handler produces a JSON rendering of a package or a dcat:CatalogRecord
// for RDF requests
func PackageService(w http.ResponseWriter, req *http.Request) {
	// xxx special case
	if req.URL.Path == "/api/rest/package/" {
		PackageList(w, req)
		return
	}
	LogRequest(req)
	resource, content_type, format := negotiate(req)
	if len(content_type) == 0 {
		unacceptable(w)
		return
	}

	slugs := strings.SplitAfterN(resource, "/package/", 2) // should not really be here...
	if len(slugs) != 2 {
		notfound(w)
		return
	}
	id := slugs[1]

	pkg, err := repo.GetPackage(id)

	if err != nil {
		notfound(w)
		return
	}

	switch {
	case format == "json":
		pkgmap := pkg.ToMap()
		buf, err := json.Marshal(pkgmap)
		if err != nil {
			servererror(w, err)
			return
		}
		w.Header().Set("Content-Type", content_type)
		content_length := fmt.Sprintf("%d", len(buf))
		w.Header().Set("Content-Length", content_length)
		if !strings.HasSuffix(req.URL.Path, "."+format) {
			w.Header().Set("Vary", "Accept")
			w.Header().Set("Content-Location", req.URL.Path+".json")
		}
		w.Header().Set("Server", server_software)
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
	case !strings.HasSuffix(req.URL.Path, "."+format):
		w.Header().Set("Vary", "Accept")
		http.Redirect(w, req, req.URL.Path+"."+format, http.StatusSeeOther)
	default: // rdf
		serializer := rdf.Serializer(format)
		defer serializer.Free()
		resource := "http://" + req.Host + "/package/" + id // xxx
		data, err := serializer.Serialize(gockan.CatalogRecord(pkg, resource), "")
		if err != nil {
			servererror(w, err)
			return
		}
		w.Header().Set("Content-Type", content_type)
		content_length := fmt.Sprintf("%d", len(data))
		w.Header().Set("Content-Length", content_length)
		w.Header().Set("Server", server_software)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
}
