package utils

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

func PagingFind(skip, limit int, m bson.M, collection string) ([]interface{}, error) {
	var list []interface{}
	query := func(c *mgo.Collection) error {
		return c.Find(m).Sort("-create_at", "-created_at").Skip(skip).Limit(limit).All(&list)
	}
	err := GetMongoCollection(collection, query)
	return list, err
}

func FindById(id interface{}, collection string) (interface{}, error) {
	var entity interface{}
	query := func(c *mgo.Collection) error {
		return c.FindId(id).One(&entity)
	}
	err := GetMongoCollection(collection, query)
	return entity, err
}

func Insert(m interface{}, collection string) error {
	query := func(c *mgo.Collection) error {
		return c.Insert(m)
	}
	err := GetMongoCollection(collection, query)
	return err
}

func FindOne(m bson.M, collection string) (interface{}, error) {
	var entity interface{}
	query := func(c *mgo.Collection) error {
		return c.Find(m).One(&entity)
	}
	err := GetMongoCollection(collection, query)
	return entity, err
}

//func Count(m bson.M, collection string) (int, error) {
//	query := func(c *mgo.Collection) (int, error) {
//		return c.Find(m).Count()
//	}
//	err := GetMongoCollection(collection, query)
//	return err
//}

func UpdateById(id bson.ObjectId, m bson.M, collection string) error {
	query := func(c *mgo.Collection) error {
		return c.UpdateId(id, m)
	}
	err := GetMongoCollection(collection, query)
	return err
}

func Update(selector bson.M, update bson.M, collection string) error {
	query := func(c *mgo.Collection) error {
		return c.Update(selector, update)
	}
	err := GetMongoCollection(collection, query)
	return err
}

func FindAll(m bson.M, collection string) ([]interface{}, error) {
	var list = make([]interface{}, 0, 100)
	query := func(c *mgo.Collection) error {
		return c.Find(m).All(&list)
	}
	err := GetMongoCollection(collection, query)
	return list, err
}

