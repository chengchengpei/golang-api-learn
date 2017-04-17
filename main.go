/*
This is for fun and learn.
Refer the following tutorials:
    https://thenewstack.io/make-a-restful-json-api-go/


Todo:
    Add DB.
    Add logic to generate ramdom code.

*/
package main

import (
    "fmt"
    "html"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "time"
)

type Link struct {
    ShortURL string
    LongURL string
    Created time.Time
}

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
    }

    return router
}

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },

    Route{
        "UrlShortener",
        "POST",
        "/",
        UrlShortener,
    },

    Route{
        "redirect",
        "POST",
        "/{code}",
        Lookup,
    },
}


func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello index, %q", html.EscapeString(r.URL.Path))
}

func UrlShortener(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Lookup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    code := vars["code"]
    fmt.Fprintln(w, "code to lookup:", code)
}


func main() {
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8000", router))
}


