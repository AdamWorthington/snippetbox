package main

import (
    "bytes"
    "fmt"
    "errors"
    "net/http"
    "time"

    "github.com/go-playground/form/v4"
)


func (app *application) decodePostForm(r *http.Request, dst any) error {
    err := r.ParseForm()
    if err != nil {
        return err
    }

    err = app.formDecoder.Decode(dst, r.PostForm)
    if err != nil {
        var invalidDecoderError *form.InvalidDecoderError
        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }
        return err
    }
    return nil
}
            

func (app *application) newTemplateData(r *http.Request) templateData {
    return templateData {
        CurrentYear: time.Now().Year(),
        Flash: app.sessionManager.PopString(r.Context(), "flash"),
    }
}
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) { // Retrieve the appropriate template set from the cache based on the page
    ts, ok := app.templateCache[page] 
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page) 
        app.serveError(w, r, err)
        return
    }
    buf := new(bytes.Buffer)

    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serveError(w, r, err) 
        return
    }
    w.WriteHeader(status)
    buf.WriteTo(w)
}
func (app *application) serveError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri = r.URL.RequestURI()
    )
    app.logger.Error(err.Error(), "Method", method, "uri", uri)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}

