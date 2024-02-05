package main

import(
    "log/slog"
    "os"
    "flag"
    "net/http"
)

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")

    flag.Parse()
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    
    mux := http.NewServeMux()
    fileserver := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileserver))

    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet/view", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)

    logger.Info("Starting server on port", "addr", *addr)

    err := http.ListenAndServe(*addr, mux)
    logger.Error(err.Error())
    os.Exit(1)
}   
