package main

import (
    "errors"
    "fmt"
    "net/http"
    "strings"
    "unicode/utf8"
    "strconv"
    
    "snippetbox.adamworthington.net/internal/models"
    "github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
    Title string
    Content string
    Expires int
    FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serveError(w, r, err)
        return
    }

    data := app.newTemplateData(r)
    data.Snippets = snippets

    app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r * http.Request) {
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
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
    data := app.newTemplateData(r)
    data.Snippet = snippet
    
    app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) { 
    data := app.newTemplateData(r)
    data.Form = snippetCreateForm {
        Expires: 365,
    }
    app.render(w, r, http.StatusOK, "create.html", data) 
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) { 
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.serveError(w, r, err)
        return
    }

    fieldErrors := make(map[string]string)

    form := snippetCreateForm {
        Title: r.PostForm.Get("title"),
        Content: r.PostForm.Get("content"),
        Expires: expires,
        FieldErrors: fieldErrors,
    }
    if strings.TrimSpace(form.Title) == "" {
        fieldErrors["title"] = "This field can not be blank"
    } else if utf8.RuneCountInString(form.Title) > 100 {
        fieldErrors["title"] = "Title can not be more than 100 characters long"
    }

    if strings.TrimSpace(form.Content) == "" {
        fieldErrors["content"] = "This field can not be blank"
    }

    if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
        fieldErrors["expires"] = "Expire date must be 1, 7, or 365 days"
    }

    if len(fieldErrors) > 0 {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
        return
    }

    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serveError(w, r, err)
        return
    }
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}   

