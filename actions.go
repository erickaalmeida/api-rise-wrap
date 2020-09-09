package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var collection = getSession().DB("rise_wrap").C("wrap")

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}
	return session
}

func ResponseWrap(w http.ResponseWriter, status int, result Wrap) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func ResponseWraps(w http.ResponseWriter, status int, result Wraps) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola mundo desde mi servidor web con Go")
}

func WrapList(w http.ResponseWriter, r *http.Request) {
	var wraps []Wrap
	err := collection.Find(nil).Sort("-_id").All(&wraps)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("resultado: ", wraps)
	}
	ResponseWraps(w, 200, wraps)
}

func WrapShow(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	wrapId := urlParams["id"]

	if !bson.IsObjectIdHex(wrapId) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(wrapId)
	var wrap Wrap
	err := collection.FindId(oid).One(&wrap)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	ResponseWrap(w, 200, wrap)
}

func WrapAdd(w http.ResponseWriter, r *http.Request) {
	wrapDecoder := json.NewDecoder(r.Body)

	var wrapData Wrap

	err := wrapDecoder.Decode(&wrapData)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	err = collection.Insert(wrapData)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	ResponseWrap(w, 200, wrapData)
}

func WrapUpdate(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	wrapId := urlParams["id"]

	if !bson.IsObjectIdHex(wrapId) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(wrapId)
	wrapDecoder := json.NewDecoder(r.Body)
	var wrapData Wrap
	errDecoder := wrapDecoder.Decode(&wrapData)

	if errDecoder != nil {
		log.Fatal(errDecoder)
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	var wrap Wrap
	err := collection.FindId(oid).One(&wrap)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	wrapDocument := bson.M{"_id": oid}
	wrapChange := bson.M{"$set": wrapData}
	errUpdate := collection.Update(wrapDocument, wrapChange)

	if errUpdate != nil {
		panic(errUpdate)
		w.WriteHeader(500)
		return
	}
	ResponseWrap(w, 200, wrapData)
}

func WrapRemove(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	wrapId := urlParams["id"]

	if !bson.IsObjectIdHex(wrapId) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(wrapId)

	err := collection.RemoveId(oid)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	removeResult := Message{
		"success",
		fmt.Sprintf("El menú con Id %v ha sido eliminado", wrapId),
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(removeResult)
}
