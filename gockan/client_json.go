package gockan

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/arjanvaneersel/gockan/model"
)

type JsonClient struct {
	ApiBase, ApiKey, ExtraBase, TagBase, GroupBase string
	client                                         http.Client
}

func NewJsonClient(args ...string) (repo Repository) {
	jc := JsonClient{}
	if len(args) > 0 {
		jc.ApiBase = args[0]
	}
	if len(args) > 1 {
		jc.ApiKey = args[1]
	}
	if len(args) > 2 {
		jc.ExtraBase = args[2]
	}
	if len(args) > 3 {
		jc.TagBase = args[3]
	}
	if len(args) > 4 {
		jc.GroupBase = args[4]
	}

	if len(jc.ApiBase) == 0 {
		jc.ApiBase = "http://ckan.net/api/"
	}
	if len(jc.ExtraBase) == 0 {
		jc.ExtraBase = "http://wiki.ckan.net/extras/"
	}
	if len(jc.TagBase) == 0 {
		jc.TagBase = "http://ckan.net/tag/"
	}
	if len(jc.GroupBase) == 0 {
		jc.GroupBase = "http://ckan.net/group/"
	}
	return &jc
}

// utility function to read the body of an HTTP response
func read_response(resp *http.Response) (body []byte, err error) {
	if resp.ContentLength >= 0 {
		body = make([]byte, 0, resp.ContentLength)
	} else {
		body = make([]byte, 0, 4096)
	}

	for {
		buf := make([]byte, 4096)
		n, errRead := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if errRead != nil {
			if errRead != io.EOF {
				err = errRead
			} else {
				err = nil
			}
			break
		}
	}
	return
}

func (jc *JsonClient) get(u string) (body []byte, err error) {
	var req http.Request
	req.Method = "GET"
	req.ProtoMajor = 1
	req.ProtoMinor = 1

	req.URL, err = url.Parse(u)
	if err != nil {
		return
	}

	log.Printf("requesting %s", req.URL)

	req.Header = make(http.Header)
	req.Header["User-Agent"] = []string{UserAgent}
	req.Header["Accept"] = []string{"application/json"}

	if len(jc.ApiKey) != 0 {
		req.Header["X-CKAN-APIKEY"] = []string{jc.ApiKey}
	}

	resp, err := jc.client.Do(&req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return
	}

	body, err = read_response(resp)
	if err != nil {
		return
	}

	return
}

func (jc *JsonClient) package_list() (pkglist []string, err error) {
	rest_url := jc.ApiBase + "rest/package"

	body, err := jc.get(rest_url)
	if err != nil {
		return
	}

	pkglist = make([]string, 0)
	err = json.Unmarshal(body, &pkglist)
	if err != nil {
		return
	}

	return
}

func (jc *JsonClient) Count() int {
	pkglist, err := jc.package_list()
	if err != nil {
		return -1
	}
	return len(pkglist)
}

func (jc *JsonClient) Packages() (ch chan *model.Package, err error) {
	rest_url := jc.ApiBase + "rest/package"

	body, err := jc.get(rest_url)
	if err != nil {
		return
	}

	packages := make([]string, len(body)/38) // list of uuids, 36 chars each + quotes
	err = json.Unmarshal(body, &packages)
	if err != nil {
		return
	}

	ch = make(chan *model.Package)
	go func() {
		for _, pkgid := range packages {
			pkg, err := jc.GetPackage(pkgid)
			if err != nil {
				log.Printf("GetPackage(%s): %s", pkgid, err)
				continue
			}
			ch <- pkg
		}
		close(ch)
	}()
	return
}

func (jc *JsonClient) GetPackage(id string) (pkg *model.Package, err error) {
	rest_url := jc.ApiBase + "rest/package/" + id
	body, err := jc.get(rest_url)
	if err != nil {
		return
	}

	pkgmap := make(map[string]interface{})
	err = json.Unmarshal(body, &pkgmap)
	if err != nil {
		return
	}

	pkg = model.PackageFromMap(pkgmap)
	pkg.Provenance = model.NewProcess()
	pkg.Provenance.SetOp(UserAgent)
	pkg.Provenance.Use(rest_url)
	pkg.ExtraBase = jc.ExtraBase
	pkg.TagBase = jc.TagBase
	pkg.GroupBase = jc.GroupBase

	return
}

func (jc *JsonClient) PutPackage(pkg *model.Package) (err error) {
	pkgmap := pkg.ToMap()

	data, err := json.Marshal(pkgmap)

	log.Printf("JSON:\n%s", string(data))
	return
}
