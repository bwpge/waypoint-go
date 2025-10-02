package main

import (
	"fmt"
	"net/http"

	log "github.com/bwpge/systemdlog"

	flag "github.com/spf13/pflag"
)

type Handler struct {
	config config
	cache  cache
}

func NewHandler() *Handler {
	h := Handler{}
	h.cache = NewCache()

	conf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	h.config = conf

	return &h
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	v, found := h.config[key]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	k := v.url
	if v.redir {
		http.Redirect(w, r, k, http.StatusTemporaryRedirect)
		return
	}

	body, found := h.cache.get(k)
	if found {
		log.Debugf("cache hit: '%s'", k)
		fmt.Fprint(w, body)
		return
	}

	log.Debugf("cache miss: '%s'", k)
	body, err := fetch(k)
	if err != nil {
		log.Errorf("failed to fetch '%s': %s", k, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.cache.set(k, body)
	fmt.Fprint(w, body)
}

func main() {
	port := flag.UintP("port", "p", 8080, "port to listen on")
	host := flag.StringP("host", "h", "", "IP address or hostname to bind")
	confPath := flag.StringP("config", "c", "", "explicit config file path instead of defaults")
	flag.Parse()

	p := *confPath
	if p != "" {
		if err := checkFile(p); err != nil {
			log.Fatalf("could not load explicit config '%s': %s", p, err)
		}
		confPaths = []string{p}
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	h := NewHandler()

	http.HandleFunc("/{key}", h.Handle)
	log.Infof("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
