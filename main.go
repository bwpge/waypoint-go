package main

import (
	"fmt"
	"net/http"

	log "github.com/bwpge/systemdlog"

	flag "github.com/spf13/pflag"
)

type configItem struct {
	value     string
	isContent bool
	redir     bool
}

type itemMap map[string]configItem

type Handler struct {
	items itemMap
	cache cache
}

func NewHandler(items map[string]yamlItem) *Handler {
	h := Handler{}
	h.cache = NewCache()
	h.items = make(itemMap)

	for k, v := range items {
		item := configItem{value: v.URL, redir: v.Redir}
		if v.Content != "" {
			item.value = v.Content
			item.isContent = true
		}

		if item.value == "" {
			log.Fatalf("error in config key '%s': a url is required", k)
		}

		for _, alias := range append(v.Alias, k) {
			if _, found := h.items[alias]; found {
				log.Fatalf("error in config key '%s': duplicate key or alias '%s'", k, alias)
			}
			h.items[alias] = item
		}
	}

	return &h
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ip := extractIP(r)
	key := r.PathValue("key")
	v, found := h.items[key]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	k := v.value
	if v.isContent {
		log.Debugf("%s - CONTENT - %s", ip, key)
		fmt.Fprint(w, v.value)
		return
	}

	if v.redir {
		log.Debugf("%s - REDIRECT - %s => %s", ip, key, k)
		http.Redirect(w, r, k, http.StatusTemporaryRedirect)
		return
	}

	body, found := h.cache.get(k)
	if found {
		log.Debugf("%s - HIT - %s => %s", ip, key, k)
		fmt.Fprint(w, body)
		return
	}
	log.Debugf("%s - MISS - %s => %s", ip, key, k)

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
	pPort := flag.UintP("port", "p", 6920, "port to listen on")
	pHost := flag.StringP(
		"host",
		"H",
		"127.0.0.1",
		"IP address or hostname to bind (use '-' for none)",
	)
	confPath := flag.StringP("config", "c", "", "explicit config file path instead of defaults")
	flag.Parse()

	p := *confPath
	if p != "" {
		if err := checkFile(p); err != nil {
			log.Fatalf("could not load explicit config '%s': %s", p, err)
		}
		confPaths = []string{p}
	}

	conf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	port := *pPort
	host := *pHost
	if port == 0 && conf.Options.Port != 0 {
		port = conf.Options.Port
	}
	if host == "" && conf.Options.Host != "" {
		host = conf.Options.Host
	}
	if host == "-" {
		host = ""
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	h := NewHandler(conf.Items)

	http.HandleFunc("/{key}", h.Handle)
	log.Infof("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
