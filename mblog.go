package main 

import (
	"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

type Note struct {
	Url string
	Title string
	Date string
	Body string
}

func loadNotes() ([]Note) {
	session, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("test").C("notes")

	notes := []Note{}
    iter := c.Find(nil).Iter()
    err = iter.All(&notes)
    if err != nil {
        panic(iter.Err())
    }

    return notes
}

func main() {
	notes := loadNotes()
    for _, note := range notes {
        fmt.Println(note.Title)
    }
}