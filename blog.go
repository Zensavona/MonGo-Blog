package main 

import (
	"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
    "net/http"
    "github.com/russross/blackfriday"
    "github.com/hoisie/mustache"
    "io/ioutil"
    "strings"
)

type Note struct {
	Url string
	Title string
	Date string
	Body string
}

var notes []Note = loadNotes()
var lenPath = len("/note/")

// less than ideal
var home, _ = ioutil.ReadFile("home.md")
var homeMarkup = string(blackfriday.MarkdownCommon(home))

func loadNotes() ([]Note) {
	session, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
    defer session.Close()
    c := session.DB("test").C("notes")
	notes := []Note{}
    iter := c.Find(nil).Iter()
    err = iter.All(&notes)
    if err != nil {
        panic(iter.Err())
    }
    return notes
}

func loadTemplate(name string) string {
    file, _ := ioutil.ReadFile(name + ".html")
    return string(file)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    var data struct {
        Notes []Note
    }
    data.Notes = notes
    rendered := mustache.RenderInLayout(homeMarkup, loadTemplate("home"), data)
    fmt.Fprintf(w, rendered)
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
    title := strings.Replace(r.URL.Path[lenPath:], ".html", "", -1)
    var rendered string
    for _, note := range notes {
        if note.Url == title {
            rendered = mustache.RenderInLayout(note.Body, loadTemplate("note"), note)
        }
    }
    if len(rendered) == 0 {
        fmt.Fprintf(w, "404 page not found")
    }
    fmt.Fprintf(w, rendered)
}

func main() {
    http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/note/", noteHandler)
    http.ListenAndServe(":8080", nil)
}