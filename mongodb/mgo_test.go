package main

import (
	"fmt"
	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
	"testing"
)

type programmer struct {
	Name  string
	Field string
	Age   byte
}

func (p programmer) string() string {
	return fmt.Sprintf("%s-%s\n", p.Name, p.Field)
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MGO(t *testing.T) {
	session, err := mgo.Dial(":27017")
	checkError(err, t)
	defer session.Close()

	collection := session.DB("programmers").C("Guys")
	if n, _ := collection.Count(); n != 0 {
		collection.RemoveAll(bson.M{})
	}
	err = collection.Insert(
		&programmer{"John Carmack", "Graphics", 1},
		&programmer{"agentzh", "Backend", 2},
		&programmer{"BradFitz", "Go", 3},
		&programmer{"need_delete", "...", 4},
	)
	checkError(err, t)

	err = collection.Update(bson.M{"name": "agentzh"}, &programmer{"Agentzh", "Backend", 2})
	checkError(err, t)
	_, err = collection.RemoveAll(bson.M{"name": "need_delete"})
	checkError(err, t)
	query := collection.Find(bson.M{"name": "need_delete"})
	if n, err := query.Count(); err == nil && n != 0 {
		t.FailNow()
	}

	guys := []programmer{}
	err = collection.Find(bson.M{"age": bson.M{"$lte": 2}}).All(&guys)
	checkError(err, t)
	for _, guy := range guys {
		fmt.Println(guy.string())
	}
}
