package main

import (
    "log" 
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("Hello from snippetbox"))
}


func snippetView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Specific content to serve"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Allow", http.MethodPost) 
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return;
    }
    w.Write([]byte("Creating a snippet"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet/view", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)

    log.Print("String server on port 4000")
    err:=http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}

