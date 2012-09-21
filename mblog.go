package main 

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Note struct {
	Url string
	Title string
	Date string
	Body string
}

func loadNotes() (Note) {
	session, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("test").C("notes")

	note := Note{}
    err = c.Find(bson.M{"_id": "some-shitty-title"}).One(&note)
    if err != nil {
            panic(err)
    }

    return note
}

func main() {
	post := loadNotes()

	fmt.Println(post.Body)
}