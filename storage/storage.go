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
		"state":    model.State,
		"lastbeat": model.Lastbeat,
	}
}

// UpsertModel save/update model to mongo
func UpsertModel(model model.Model, db *mgo.Database) (interface{}, error) {
	var id int32
	c := db.C("models")
	if model.ID != 0 {
		id = model.ID
	} else {
		id = bson.NewObjectId().Counter()
	}
	info, err := c.UpsertId(id, createModelBson(model))
	if err != nil {
		return "", err
	}
	fmt.Println(model.ID)
	fmt.Println(info.UpsertedId)
	fmt.Println(info.Updated)
	if info.UpsertedId == nil {
		return model.ID, nil
	}
	return info.UpsertedId, nil
}

// FindModels use query to find models obj
func FindModels(query map[string]interface{}, id interface{}, limit int, db *mgo.Database) ([]model.Model, error) {
	c := db.C("models")
	var result []bson.M
	var models []model.Model

	if id != nil {
		println("find id", id.(int32))
		err := c.FindId(id.(int32)).All(&result)
		if err != nil {
			return models, err
		}
	} else {
		err := c.Find(query).Limit(limit).All(&result)
		if err != nil {
			return models, err
		}
	}

	fmt.Println("result", result)

	for _, v := range result {
		models = append(models, model.Model{
			ID:        int32(v["_id"].(int)),
			ModelType: model.ModelType{v["type"].(string)},
			Address:   v["address"].(string),
			State:     model.ModelState(uint8(v["state"].(int))),
			Lastbeat:  int32(v["lastbeat"].(int)),
		})
	}

	return models, nil
}

// MongoErrorHandle handle redis error
func MongoErrorHandle(err error, exit bool, msg string) {
	fmt.Printf("Mongo %s\nError: %s", msg, err)
	if exit {
		panic(err)
	}
}
