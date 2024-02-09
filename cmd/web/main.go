package main

import(
    "log/slog"
    "database/sql"
    "os"
    "flag"
    "net/http"
    "snippetbox.adamworthington.net/internal/models"
    _ "github.com/go-sql-driver/mysql"
)

type application struct {
    logger *slog.Logger
    snippets *models.SnippetModel
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")

    dsn := flag.String("dsn", "web@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()
    
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close() 
    
    app := &application{
        snippets: &models.SnippetModel{DB: db},
        logger: logger,
    }

    logger.Info("Starting server on port", "addr", *addr)

    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}   

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}
