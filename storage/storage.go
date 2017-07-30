package storage

import (
	"fmt"
	"link/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// InitDatabase init mongo db with inkdb
func InitDatabase(url string) *mgo.Database {
	session, err := mgo.Dial("localhost")
	if err != nil {
		MongoErrorHandle(err, true, "can not create mongodb session")
	}

	return session.DB("linkdb")
}

func createModelBson(model model.Model) *bson.M {
	return &bson.M{
		"type":     model.ModelType.Name,
		"address":  model.Address,
		"state":    string(model.State),
		"lastbeat": string(model.Lastbeat),
	}
}

// UpsertModel save/update model to mongo
func UpsertModel(model model.Model, db *mgo.Database) (interface{}, error) {
	c := db.C("models")
	info, err := c.UpsertId(model.ID, createModelBson(model))
	if err != nil {
		return "", err
	}

	return info.UpsertedId, nil
}

// MongoErrorHandle handle redis error
func MongoErrorHandle(err error, exit bool, msg string) {
	fmt.Printf("Mongo %s\nError: %s", msg, err)
	if exit {
		panic(err)
	}
}
