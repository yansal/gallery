package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yansal/gallery/storage"
	"github.com/yansal/gallery/storage/s3"
)

func main() {
	bucket := os.Getenv("S3BUCKET")
	if bucket == "" {
		log.Fatal("S3BUCKET is required")
	}
	storage, err := s3.New(bucket)
	if err != nil {
		log.Fatal(err)
	}

	templates := template.Must(template.ParseGlob("templates/*"))
	http.Handle("/", &handler{
		storage:   storage,
		templates: templates,
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type handler struct {
	storage   storage.Storage
	templates *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}

	herr, ok := err.(httpError)
	if !ok {
		herr = httpError{err: err, code: http.StatusInternalServerError}
	}
	http.Error(w, fmt.Sprintf("%+v", herr.Error()), herr.code)
}

type httpError struct {
	err  error
	code int
}

func (e httpError) Error() string { return e.err.Error() }

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/" {
		res, err := h.storage.List("")
		if err != nil {
			return err
		}
		if err := h.templates.ExecuteTemplate(w, "index.html", res); err != nil {
			log.Print(err)
		}
	} else {
		res, err := h.storage.List(r.URL.Path[1:])
		if err != nil {
			return err
		}
		data := struct {
			ListResult storage.ListResult
			ImgAddr    string
		}{
			ListResult: res,
			ImgAddr:    os.Getenv("IMGADDR"),
		}
		if err := h.templates.ExecuteTemplate(w, "gallery.html", data); err != nil {
			log.Print(err)
		}
	}
	return nil
}
