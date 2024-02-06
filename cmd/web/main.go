package main

import(
    "log/slog"
    "os"
    "flag"
    "net/http"
)

type application struct {
    logger *slog.Logger
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()
    
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    
    app := &application{
        logger: logger,
    }
    mux := http.NewServeMux()
    fileserver := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileserver))

    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snipper/view", app.snippetView)
    mux.HandleFunc("/snippet/create", app.snippetCreate)

    logger.Info("Starting server on port", "addr", *addr)

    err := http.ListenAndServe(*addr, mux)
    logger.Error(err.Error())
    os.Exit(1)
}   
