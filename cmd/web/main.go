package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"snippetbox.adamworthington.net/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
    logger *slog.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
    formDecoder *form.Decoder
    sessionManager *scs.SessionManager
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
    cache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    
    formDecoder := form.NewDecoder()

    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = time.Hour * 12
    sessionManager.Cookie.Secure = true

    app := &application{
        snippets: &models.SnippetModel{DB: db},
        logger: logger,
        templateCache : cache,
        formDecoder: formDecoder,
        sessionManager: sessionManager,
    }

    srv := &http.Server{
        Addr: *addr,
        Handler: app.routes(),
        ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }

    logger.Info("Starting server on port", "addr", *addr)

    err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
