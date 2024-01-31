package main

import (
    "fmt"
    "strconv"
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
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 0 {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "Display a specific snippet with id: %d...", id)
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

