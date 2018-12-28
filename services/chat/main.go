package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once // compile the template once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP Request, by handling the loading, compiling and delivering the template
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "The addr of the applicaiton")
	//  parse the flags
	flag.Parse()
	r := newRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	fmt.Printf("Listening on port %v...\n", *addr)
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
