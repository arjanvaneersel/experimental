package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/ww/goraptor"
	"github.com/vuleetu/goconfig/config"
	"github.com/arjanvaneersel/gockan"
)

var repo gockan.Repository

var daemonise bool
var cfgfile, cwd, program string
var socket *net.TCPListener

const server_version = "0.1.1"
const server_software = "ckand v0.1.1"

// holds the list of content-types that we can handle
var mime_types []string

// holds a map content_type -> raptor serializer
var mime_map map[string]string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", server_software)
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] config\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.BoolVar(&daemonise, "d", false, "Daemonise")

	// for our persistence, gob needs to be told about these types
	uri := goraptor.Uri("")
	gob.Register(&uri)
	bnode := goraptor.Blank("")
	gob.Register(&bnode)
	literal := goraptor.Literal{}
	gob.Register(&literal)

	// set up mime types <-> format mappings
	mime_types = make([]string, 0, len(goraptor.SerializerSyntax)+1)
	mime_map = make(map[string]string)

	mime_map["text/javascript"] = "json"
	mime_types = append(mime_types, "text/javascript")
	mime_map["application/json"] = "json"
	mime_types = append(mime_types, "application/json")

	for _, syntax := range goraptor.SerializerSyntax {
		// do not step on our own native json format
		if syntax.MimeType == "application/json" {
			continue
		}
		mime_map[syntax.MimeType] = syntax.Name
		mime_types = append(mime_types, syntax.MimeType)
	}

}

func main() {
	var err error

	// parse flags
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(255)
	}

	// set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logpfx := fmt.Sprintf("[ckand %d] ", os.Getpid())
	log.SetPrefix(logpfx)

	// find and save the current working directory
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// get the full path to the program
	program = os.Args[0]
	if !filepath.IsAbs(program) {
		if strings.Index(program, "/") >= 0 {
			program = filepath.Join(cwd, program)
			program = filepath.Clean(program)
		} else {
			program, err = exec.LookPath(program)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// get the absolute path to the configuration file
	cfgfile = flag.Arg(0)
	if !filepath.IsAbs(cfgfile) {
		cfgfile = filepath.Join(cwd, cfgfile)
		cfgfile = filepath.Clean(cfgfile)
	}

	// read the configuration file
	cfg, err := config.ReadDefault(cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	// daemonise if we have been told to
	if daemonise {
		daemon()
	}

	// change to the configured database root if necessary
	root, err := cfg.String("DEFAULT", "root")
	log.Print("root ", root)
	if err == nil {
		if !filepath.IsAbs(root) {
			root = filepath.Join(cwd, root)
			root = filepath.Clean(root)
		}
		err = os.MkdirAll(root, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(root)
		if err != nil {
			log.Fatal(err)
		}
	}

	// start logging to file if necessary
	logfile, err := cfg.String("log", "file")
	if err == nil {
		fp, err := os.Open(logfile) //, os.O_CREAT|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(fp)
	}

	log.Printf("starting up with config %s", cfgfile)

	// set maxprocs
	maxprocs, err := cfg.Int("DEFAULT", "maxprocs")
	if err == nil {
		runtime.GOMAXPROCS(maxprocs)
	}
	log.Printf("gomaxprocs: %d", runtime.GOMAXPROCS(0))

	// listen on the appropriate address and port
	addr, err := cfg.String("http", "bind")
	if err != nil {
		addr = ":8080"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	socket = listener.(*net.TCPListener)
	log.Printf("listening on %s", addr)

	// create our in memory repository
	repo = gockan.NewMemRepo()

	// create data sources to aggregate
	err = init_aggregator(cfg)
	if err != nil {
		log.Fatal(err)
	}

	handle_cneg("/catalogue", PackageList)
	http.HandleFunc("/package/", PackageService)

	// ckan api compat
	handle_cneg("/api/rest/package", PackageList)
	http.HandleFunc("/api/rest/package/", PackageService)
	handle_cneg("/api/1/rest/package", PackageList)
	http.HandleFunc("/api/1/rest/package/", PackageService)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req,
			"http://godoc.styx.org/pkg/bitbucket.org/okfn/gockan/ckand",
			http.StatusMovedPermanently)
	})

	go SignalHandler(cfg)

	http.Serve(listener, nil)
}

func handle_cneg(base string, handler http.HandlerFunc) {
	http.HandleFunc(base, handler)
	for _, syntax := range goraptor.SerializerSyntax {
		http.HandleFunc(base+"."+syntax.Name, handler)
	}
}

func init_aggregator(cfg *config.Config) (err error) {
	sources, err := cfg.String("aggregator", "sources")
	if err != nil {
		err = nil // aggregation is optional
		return
	}

	for _, source := range strings.Fields(sources) {
		disabled, _ := cfg.Bool(source, "disabled")
		if disabled {
			continue
		}

		log.Print("configuring source ", source)
		var src gockan.Repository

		source_type, _ := cfg.String(source, "type")
		extra_base, _ := cfg.String(source, "extra_base")
		group_base, _ := cfg.String(source, "group_base")
		tag_base, _ := cfg.String(source, "tag_base")
		switch source_type {
		case "json":
			api_base, _ := cfg.String(source, "api_base")
			api_key, _ := cfg.String(source, "api_key")
			_, _ = LoadRepo(repo, source+".gob")
			src = gockan.NewJsonClient(api_base, api_key, extra_base, tag_base, group_base)
		default:
			errs := fmt.Sprintf("unknown source type [%s]: %s", source, source_type)
			err = errors.New(errs)
			return
		}
		harvest, _ := cfg.Bool(source, "harvest")
		if harvest {
			harvester := NewHarvester(cfg, source, src, repo)
			harvester.Start()
		}
	}
	return
}
