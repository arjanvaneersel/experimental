package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"github.com/arjanvaneersel/go-chat/trace"
	"os"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var avatars Avatar = TryAvatars{
	UseAuthAvatar,
	UseGravatar,
	UseFileSystemAvatar,
}

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "Application's address")
	flag.Parse()

	gomniauth.SetSecurityKey("ejtr43u3j0ty8uj53ujqjOIETOIGHI456JY4#it$")
	gomniauth.WithProviders(
		google.New(
			"767894946599-cm9o5qrtddnl24nm8k2sqbkgvrapm7f0.apps.googleusercontent.com",
			"X-EY97dEiLU_e5_8isI0J6KB",
			"http://localhost:8080/auth/callback/google",
		))

	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/upload", &templateHandler{filename:"upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.Handle("/assets", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	go r.run()

	log.Println("Starting server on ", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
