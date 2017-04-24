package main

import (
    "fmt"
    "html"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    // "time"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func randCode() string {
    b := make([]rune, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

type Link struct {
    ShortURL string
    LongURL string
    // Created time.Time
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
        "GET",
        "/{code}",
        Lookup,
    },
}


func Index(w http.ResponseWriter, r *http.Request) {
    // Set status code explicitly
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello index, %q", html.EscapeString(r.URL.Path))
}

func UrlShortener(w http.ResponseWriter, r *http.Request) {
    r.ParseForm() // Parse the request body
    longUrl := r.Form.Get("longUrl") // longUrl will be '' if no longUrl in the requst

    /* DB operations */
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    c := session.DB("test").C("link")
    result := Link{}
    err = c.Find(bson.M{"longurl": longUrl}).One(&result)
    if err == nil {
        fmt.Println("Found?")
        if err := json.NewEncoder(w).Encode(result); err != nil {
            panic(err)
        }
        return
    }
    code := "" 
    fmt.Println("generating code...")

    for {
        // Find
        code = randCode()
        result := Link{}
        err = c.Find(bson.M{"shorturl": code}).One(&result)
        if err != nil && err == mgo.ErrNotFound {
            fmt.Println(err)
            err2 := c.Insert(&Link{code, longUrl})
            if err2 != nil {
                fmt.Println(err2)
            } else {
                break
            }
        }
    }

    obj :=  Link{
        LongURL: longUrl,
        ShortURL: code,
    }

    if err := json.NewEncoder(w).Encode(obj); err != nil {
        panic(err)
    }
}

func Lookup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    code := vars["code"]

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    c := session.DB("test").C("link")
    result := Link{}
    err = c.Find(bson.M{"shorturl": code}).One(&result)
    if err != nil && err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "code not found:", code)
    } else {
		http.Redirect(w, r, result.LongURL, 301)
	}

}

func main() {
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8000", router))
}


