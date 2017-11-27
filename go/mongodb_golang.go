package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	id 	string  `bson:"id"`
	Name       string    `bson:"name"`
	Lastname string      `bson:"lastname"`
	Age     int    `bson:"age"`
	Email string `bson:"email"`
	Rg      int `bson:"rg"`
	Cpf     int  `bson:"cpf"`
	Cnpj	int  `bson:"cnpj"`
	Street 	string `bson:"street"`
	Neighbourhood string `bson:"Neighbourhood"`
	City string `bson:"city"`
	State string `bson:"state"`

}

func main() {  
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    ensureIndex(session)

    mux := goji.NewMux()
    mux.HandleFunc(pat.Post("/person/"), addPerson(session))
    mux.HandleFunc(pat.Put("/person/:id"), updatePerson(session))
    mux.HandleFunc(pat.Delete("/person/:id"), deletePerson(session))
    http.ListenAndServe("localhost:8080", mux)
}

func addPerson(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {  
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        var person Person
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&person)
        if err != nil {
            ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("cadastro").C("person")

        err = c.Insert(person)
        if err != nil {
            if mgo.IsDup(err) {
                ErrorWithJSON(w, "Pessoa já cadastrada", http.StatusBadRequest)
                return
            }

            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Erro ao cadastrar usuário: ", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Location", r.URL.Path+"/"+person.id)
        w.WriteHeader(http.StatusCreated)
    }
}

func deletePerson(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {  
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        id := pat.Param(r, "id")

        c := session.DB("cadastro").C("person")

        err := c.Remove(bson.M{"id": id})
        if err != nil {
            switch err {
            default:
                ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Erro ao deletar usuário: ", err)
                return
            case mgo.ErrNotFound:
                ErrorWithJSON(w, "Usuário não encontrado", http.StatusNotFound)
                return
            }
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func updatePerson(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {  
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        id := pat.Param(r, "id")

        var person Person
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&person)
        if err != nil {
            ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("cadastro").C("person")

        err = c.Update(bson.M{"id": id}, &person)
        if err != nil {
            switch err {
            default:
                ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Erro ao atualizar o usuário: ", err)
                return
            case mgo.ErrNotFound:
                ErrorWithJSON(w, "Usuário não encontrado", http.StatusNotFound)
                return
            }
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
