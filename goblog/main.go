package goblog

import (
	"net/http"
	"github.com/gorilla/mux"
	"flag"
	"google.golang.org/appengine/user"
)

func init() {
	staticDir := flag.String("staticdir", "./static", "Path to the static directory")
	staticUrl := flag.String("staticurl", "/static/", "Static URL")
	flag.Parse()

	DefaultTemplates = append(DefaultTemplates, "jumbotron.html")

	router := mux.NewRouter()
	router.PathPrefix(*staticUrl).Handler(http.StripPrefix(*staticUrl, http.FileServer(http.Dir(*staticDir))))

	homeView, _ := NewView("base", "GoBlog", "home.html")
	router.Handle("/", homeView)

	router.HandleFunc("/{slug}", ArticleHandler)
	router.HandleFunc("/{lang:[a-z]{2}}/{slug}", ArticleHandler)

	http.Handle("/", router)
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	lang := vars["lang"]

	view, _ := NewView("base", "Article", "article.html")
	view.Data["Lang"] = lang
	view.Data["Slug"] = slug
	view.Execute(w)
}