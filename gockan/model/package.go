package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arjanvaneersel/gockan/mime"
	"github.com/arjanvaneersel/gockan/rdf"
	"bitbucket.org/ww/goraptor"
)

// An agent that typically has some relationship to a
// package. This can be a maintainer or an author and
// is the homologue of a foaf:Agent
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

func (agent *Agent) ToRdf() (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	switch agent.Node.(type) {
	case *goraptor.Blank:
		agent.Node = rdf.Blank()
	case nil:
		agent.Node = rdf.Blank()
	}

	go func() {
		if len(agent.Name) != 0 {
			name := goraptor.Literal{Value: agent.Name}
			ch <- &goraptor.Statement{agent.Node, rdf.FOAF.U("name"), &name, nil}
		}
		if len(agent.Mbox) != 0 {
			mbox := goraptor.Uri("mailto:" + agent.Mbox)
			ch <- &goraptor.Statement{agent.Node, rdf.FOAF.U("mbox"), &mbox, nil}
		}
		if len(agent.Homepage) != 0 {
			homepage := goraptor.Uri(agent.Homepage)
			ch <- &goraptor.Statement{agent.Node, rdf.FOAF.U("homepage"), &homepage, nil}
		}
		close(ch)
	}()

	return
}

// Resource mirroring package resources in CKAN
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

func (res *Resource) ToRdf() (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	switch res.Node.(type) {
	case *goraptor.Blank:
		res.Node = rdf.Blank()
	case nil:
		res.Node = rdf.Blank()
	}
	go func() {
		ch <- &goraptor.Statement{res.Node, rdf.RDF.U("type"), rdf.DCAT.U("Distribution"), nil}
		if len(res.Format) != 0 {
			format := rdf.Blank()
			ch <- &goraptor.Statement{res.Node, rdf.DC.U("format"), format, nil}
			ch <- &goraptor.Statement{format, rdf.RDF.U("type"), rdf.DC.U("IMT"), nil}
			mt, ok := mime.MimeTypes[res.Format]
			tag := rdf.Blank()
			ch <- &goraptor.Statement{tag, rdf.RDF.U("type"), rdf.MOAT.U("Tag"), nil}
			name := goraptor.Literal{Value: res.Format}
			ch <- &goraptor.Statement{tag, rdf.MOAT.U("name"), &name, nil}
			ch <- &goraptor.Statement{format, rdf.MOAT.U("taggedWithTag"), tag, nil}
			if ok {
				label := goraptor.Literal{Value: mt.Label}
				ch <- &goraptor.Statement{format, rdf.RDFS.U("label"), &label, nil}
				value := goraptor.Literal{Value: mt.Value}
				ch <- &goraptor.Statement{format, rdf.RDF.U("value"), &value, nil}
			}
		}
		if len(res.AccessURL) != 0 {
			url := goraptor.Uri(res.AccessURL)
			ch <- &goraptor.Statement{res.Node, rdf.DCAT.U("accessURL"), &url, nil}
		}
		if len(res.Description) != 0 {
			desc := goraptor.Literal{Value: res.Description}
			ch <- &goraptor.Statement{res.Node, rdf.DC.U("description"), &desc, nil}
		}
		close(ch)
	}()
	return ch
}

// Create a resource from a map representation. This will
// typically be called with the results of json.Unmarshal
func ResourceFromMap(resmap map[string]interface{}) (res *Resource) {
	res = &Resource{}
	switch id := resmap["id"].(type) {
	case string:
		res.Id = id
	}
	switch url := resmap["url"].(type) {
	case string:
		res.AccessURL = url
	}
	switch pkgid := resmap["package_id"].(type) {
	case string:
		res.Package = pkgid
	}
	switch group := resmap["resource_group_id"].(type) {
	case string:
		res.Group = group
	}
	switch format := resmap["format"].(type) {
	case string:
		res.Format = format
	}
	switch descr := resmap["description"].(type) {
	case string:
		res.Description = descr
	}
	switch hash := resmap["hash"].(type) {
	case string:
		res.Hash = hash
	}
	switch pos := resmap["position"].(type) {
	case float64:
		res.Position = int(pos)
	}
	return
}

// Return a map representation of the resource suitable for
// passing to json.Marshal
func (res *Resource) ToMap() (resmap map[string]interface{}) {
	resmap = make(map[string]interface{})
	if len(res.Id) != 0 {
		resmap["resource_id"] = res.Id
	}
	if len(res.Package) != 0 {
		resmap["package_id"] = res.Package
	}
	if len(res.Group) != 0 {
		resmap["resource_group_id"] = res.Group
	}
	if len(res.AccessURL) != 0 {
		resmap["url"] = res.AccessURL
	}
	if len(res.Format) != 0 {
		resmap["format"] = res.Format
	}
	if len(res.Hash) != 0 {
		resmap["hash"] = res.Hash
	}
	if len(res.Description) != 0 {
		resmap["description"] = res.Description
	}
	resmap["position"] = res.Position

	return
}

// Package structure that mirrors that present in CKAN
type Package struct {
	Node                          goraptor.Term
	ExtraBase, TagBase, GroupBase string
	Provenance                    *Process
	Id                            string
	Revision                      string
	State                         string
	Created                       time.Time
	Modified                      time.Time
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

// Create a package from a map. This will typically be called with the
// results of json.Unmarshal
func PackageFromMap(pkgmap map[string]interface{}) (pkg *Package) {
	pkg = &Package{}

	id := pkgmap["id"]
	switch v := id.(type) {
	case string:
		pkg.Id = v
	}

	rev := pkgmap["revision_id"]
	switch v := rev.(type) {
	case string:
		pkg.Revision = v
	}

	state := pkgmap["state"]
	switch v := state.(type) {
	case string:
		pkg.State = v
	}

	created := pkgmap["metadata_created"]
	switch v := created.(type) {
	case string:
		var err error
		tsp := strings.SplitAfterN(v, ".", 2)
		pkg.Created, err = time.Parse(time.RFC3339, tsp[0]+"Z")
		if err != nil {
			log.Print("could not parse created: ", err)
		}
	}

	modified := pkgmap["metadata_modified"]
	switch v := modified.(type) {
	case string:
		var err error
		tsp := strings.SplitAfterN(v, ".", 2)
		pkg.Modified, err = time.Parse(time.RFC3339, tsp[0]+"Z")
		if err != nil {
			log.Print("could not parse modified: ", err)
		}
	}

	url := pkgmap["ckan_url"]
	switch v := url.(type) {
	case string:
		uri := goraptor.Uri(v)
		pkg.Node = &uri
	}

	homepage := pkgmap["url"]
	switch v := homepage.(type) {
	case string:
		uri := goraptor.Uri(v)
		pkg.Homepage = &uri
	}

	cname := pkgmap["name"]
	switch v := cname.(type) {
	case string:
		pkg.Name = v
	}

	title := pkgmap["title"]
	switch v := title.(type) {
	case string:
		pkg.Title = v
	}

	version := pkgmap["version"]
	switch v := version.(type) {
	case string:
		pkg.Version = v
	}

	author := pkgmap["author"]
	author_mbox := pkgmap["author_email"]
	if author != nil || author_mbox != nil {
		pkg.Author = &Agent{}
	}
	switch v := author.(type) {
	case string:
		pkg.Author.Name = v
	}
	switch v := author_mbox.(type) {
	case string:
		pkg.Author.Mbox = v
	}

	maint := pkgmap["maintainer"]
	maint_mbox := pkgmap["maintainer_email"]
	if maint != nil || maint_mbox != nil {
		pkg.Maintainer = &Agent{}
	}
	switch v := maint.(type) {
	case string:
		pkg.Maintainer.Name = v
	}
	switch v := maint_mbox.(type) {
	case string:
		pkg.Maintainer.Mbox = v
	}

	license_id := pkgmap["license_id"]
	license := pkgmap["license"]
	if license_id != nil && license != nil {
		switch lid := license_id.(type) {
		case string:
			pkg.License = lid
			_, ok := Licenses[lid]
			if !ok {
				switch lic := license.(type) {
				case string:
					Licenses[lid] = lic
				}
			}
		}
	}

	tags := pkgmap["tags"]
	if tags != nil {
		switch tts := tags.(type) {
		case []interface{}:
			pkg.Tags = make([]string, len(tts))
			for i, v := range tts {
				pkg.Tags[i] = fmt.Sprintf("%v", v)
			}
		}
	}

	groups := pkgmap["groups"]
	if groups != nil {
		switch grs := groups.(type) {
		case []interface{}:
			pkg.Groups = make([]string, len(grs))
			for i, v := range grs {
				pkg.Groups[i] = fmt.Sprintf("%v", v)
			}
		}
	}

	extras := pkgmap["extras"]
	switch e := extras.(type) {
	case map[string]interface{}:
		pkg.Extras = make(map[string]string)
		for k, v := range e {
			pkg.Extras[k] = fmt.Sprintf("%v", v)
		}
	}

	resources := pkgmap["resources"]
	if resources != nil {
		switch rr := resources.(type) {
		case []interface{}:
			pkg.Resources = make([]*Resource, len(rr))
			for i, vi := range rr {
				switch v := vi.(type) {
				case map[string]interface{}:
					pkg.Resources[i] = ResourceFromMap(v)
				}
			}
		}
	}

	ratings_count := pkgmap["ratings_count"]
	switch v := ratings_count.(type) {
	case float64:
		pkg.RatingsCount = int(v)
	}
	ratings_average := pkgmap["ratings_average"]
	switch v := ratings_average.(type) {
	case float64:
		pkg.RatingsAvg = v
	}

	notes := pkgmap["notes"]
	switch v := notes.(type) {
	case string:
		pkg.Notes = v
	}

	notes_rendered := pkgmap["notes_rendered"]
	switch v := notes_rendered.(type) {
	case string:
		pkg.NotesRendered = v
	}

	/*
		relationships := pkgmap["relationships"]
		switch v := relationships.(type) {
		case []interface{}:
			pkg.Relationships = v
		}
	*/
	// pkgmap["download_url"] = nil, false // XXX should be in resources
	return
}

func (pkg *Package) ToMap() (pkgmap map[string]interface{}) {
	pkgmap = make(map[string]interface{})

	if len(pkg.Id) != 0 {
		pkgmap["id"] = pkg.Id
	}

	if len(pkg.Revision) != 0 {
		pkgmap["revision_id"] = pkg.Revision
	}

	if len(pkg.State) != 0 {
		pkgmap["state"] = pkg.State
	}

	pkgmap["metadata_created"] = pkg.Created.Format(time.RFC3339)

	pkgmap["metadata_modified"] = pkg.Modified.Format(time.RFC3339)

	if pkg.Node != nil {
		switch uri := pkg.Node.(type) {
		case *goraptor.Uri:
			pkgmap["ckan_url"] = string(*uri)
		}
	}

	if pkg.Homepage != nil {
		switch uri := pkg.Homepage.(type) {
		case *goraptor.Uri:
			pkgmap["url"] = string(*uri)
		}
	}

	if len(pkg.Name) != 0 {
		pkgmap["name"] = pkg.Name
	}

	if len(pkg.Title) != 0 {
		pkgmap["title"] = pkg.Title
	}

	if len(pkg.Version) != 0 {
		pkgmap["version"] = pkg.Version
	}

	if pkg.Author != nil {
		if len(pkg.Author.Name) != 0 {
			pkgmap["author"] = pkg.Author.Name
		}
		if len(pkg.Author.Mbox) != 0 {
			pkgmap["author_email"] = pkg.Author.Mbox
		}
	}

	if pkg.Maintainer != nil {
		if len(pkg.Maintainer.Name) != 0 {
			pkgmap["maintainer"] = pkg.Maintainer.Name
		}
		if len(pkg.Maintainer.Mbox) != 0 {
			pkgmap["maintainer_email"] = pkg.Maintainer.Mbox
		}
	}

	if len(pkg.License) != 0 {
		pkgmap["license_id"] = pkg.License
		license, ok := Licenses[pkg.License]
		if ok {
			pkgmap["license"] = license
		}
	}

	if pkg.Tags != nil {
		pkgmap["tags"] = pkg.Tags
	}

	if pkg.Groups != nil {
		pkgmap["groups"] = pkg.Groups
	}

	if pkg.Extras != nil {
		pkgmap["extras"] = pkg.Extras
	}

	if pkg.Resources != nil {
		resources := make([]map[string]interface{}, len(pkg.Resources))
		for i, r := range pkg.Resources {
			resources[i] = r.ToMap()
		}
		pkgmap["resources"] = resources
	}

	if pkg.RatingsCount > 0 {
		pkgmap["ratings_count"] = pkg.RatingsCount
		pkgmap["ratings_average"] = pkg.RatingsAvg
	}

	if len(pkg.Notes) != 0 {
		pkgmap["notes"] = pkg.Notes
	}
	if len(pkg.NotesRendered) != 0 {
		pkgmap["notes_rendered"] = pkg.NotesRendered
	}

	if pkg.Relationships != nil {
		pkgmap["relationships"] = pkg.Relationships
	}

	return
}

func (pkg *Package) ToRdf() (ch chan *goraptor.Statement) {
	ch = make(chan *goraptor.Statement)
	switch pkg.Node.(type) {
	case *goraptor.Blank:
		pkg.Node = rdf.Blank()
	case nil:
		pkg.Node = rdf.Blank()
	}
	go func() {
		ch <- &goraptor.Statement{pkg.Node, rdf.RDF.U("type"), rdf.DCAT.U("Dataset"), nil}

		name := goraptor.Literal{Value: pkg.Name}
		ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("identifier"), &name, nil}

		if len(pkg.Id) != 0 {
			id := goraptor.Literal{Value: pkg.Id}
			ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("identifier"), &id, nil}
		}

		if len(pkg.Title) != 0 {
			title := goraptor.Literal{Value: pkg.Title}
			ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("title"), &title, nil}
		}

		if pkg.Author != nil {
			author := pkg.Author.ToRdf()
			for {
				statement, ok := <-author
				if !ok {
					break
				}
				ch <- statement
			}
			ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("contributor"), pkg.Author.Node, nil}
		}

		if pkg.Maintainer != nil {
			maint := pkg.Maintainer.ToRdf()
			for {
				statement, ok := <-maint
				if !ok {
					break
				}
				ch <- statement
			}
			ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("maintainer"), pkg.Maintainer.Node, nil}
		}

		if len(pkg.License) > 0 {
			ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("rights"), rdf.LIC.U(pkg.License), nil}
		}

		if pkg.Homepage != nil {
			ch <- &goraptor.Statement{pkg.Node, rdf.FOAF.U("homepage"), pkg.Homepage, nil}
		}

		if pkg.Groups != nil {
			if len(pkg.GroupBase) == 0 {
				pkg.GroupBase = "http://ckan.net/group/"
			}
			for _, group := range pkg.Groups {
				group := goraptor.Uri(pkg.GroupBase + group)
				ch <- &goraptor.Statement{pkg.Node, rdf.DC.U("isPartOf"), &group, nil}
			}
		}

		if pkg.Tags != nil {
			if len(pkg.TagBase) == 0 {
				pkg.TagBase = "http://ckan.net/tag/"
			}
			for _, v := range pkg.Tags {
				tag := goraptor.Uri(pkg.TagBase + v)
				ch <- &goraptor.Statement{pkg.Node, rdf.MOAT.U("taggedWithTag"), &tag, nil}
			}
		}

		if pkg.Extras != nil {
			if len(pkg.ExtraBase) == 0 {
				pkg.ExtraBase = "http://wiki.ckan.net/extras/"
			}
			for k, v := range pkg.Extras {
				tag_uri := goraptor.Uri(pkg.ExtraBase + k)
				tag := &tag_uri
				value := goraptor.Literal{Value: v}
				ch <- &goraptor.Statement{pkg.Node, tag, &value, nil}
			}
		}

		if pkg.Resources != nil {
			for _, res := range pkg.Resources {
				statements := res.ToRdf()
				for {
					statement, ok := <-statements
					if !ok {
						break
					}
					ch <- statement
				}
				ch <- &goraptor.Statement{pkg.Node, rdf.DCAT.U("distribution"), res.Node, nil}
			}
		}

		if pkg.RatingsCount > 0 {
			rating := goraptor.Literal{}
			rating.Value = fmt.Sprintf("%.02f", pkg.RatingsAvg)
			rating.Datatype = rdf.XSD.U("float").String()
			ch <- &goraptor.Statement{pkg.Node, rdf.REV.U("rating"), &rating, nil}
		}

		close(ch)
	}()
	return ch
}
