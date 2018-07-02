package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/**
* TODO
* . write cache layer
**/

// - MARK: Handlers section.

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("(indexHandler) invoked.")
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("(mainHandler) invoked.")
}

func (s *Server) tokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("(tokenHandler) invoked.")
}

func (s *Server) statHandler(w http.ResponseWriter, r *http.Request) {
	var (
		output string
		err    error
	)
	output, err = s.HTTPServer.Varz.ToJSON()
	if err != nil {
		log.Println("(statHandler) unable to create json repr.")
		return
	}
	log.Println("(statHandler) invoked.")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	var (
		components []string
		content    []byte
		err        error
		uri        string // URI head component
		fpath      string // filepath
	)
	components = ParseURI(r.URL.Path)
	if len(components) <= 1 {
		log.Println("(fileHandler) unable to parse uri.", components, len(components), components[2])
		return
	}
	uri = components[2]
	log.Println("(fileHandler) invoked. URI ->", uri, " original: ", r.URL.Path)
	w.Header().Set("Content-Type", "application/binary")
	fpath = strings.Join([]string{"/tmp", uri}, "/")
	fmt.Println("fpath:", fpath)
	content, err = readFile(fpath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("(fileHanler) unable to read file. Got following error:", err)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println("(fileHandler) unable to write response data.")
		return
	}
	w.WriteHeader(http.StatusOK)
	s.mu.Lock()
	s.Counter.Add(uri)
	s.mu.Unlock()
	return
}
