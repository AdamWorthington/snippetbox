package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"
    
    "snippetbox.adamworthington.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
    snippets, err := app.snippets.Latest()

    if err != nil {
        app.serveError(w, r, err)
        return
    }
    app.render(w, r, http.StatusOK, "home.html", templateData{
        Snippets: snippets,
    })
}

func (app *application) snippetView(w http.ResponseWriter, r * http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serveError(w, r, err)
        }
        return
    }
    app.render(w, r, http.StatusOK, "view.html", templateData{
        Snippet: snippet,
    })
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) { 
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }
    title := "0 snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa" 
    expires := 7

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serveError(w, r, err)
        return
    }
    http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}   

