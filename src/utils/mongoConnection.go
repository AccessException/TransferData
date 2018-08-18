package utils

import (
	"gopkg.in/mgo.v2"
	"log"
)

// mongo会话
var mgoSession *mgo.Session

// 获取mongo会话
func getSession() *mgo.Session {
	if mgoSession == nil {
		// 错误信息
		var err error
		// 连接mongo数据库
		mgoSession, err = mgo.Dial(Config.MongoAddress)
		// 连接数据库异常
		if err != nil {
			log.Println("连接mongo数据库失败！")
			//直接终止程序运行
			panic(err)
		}else{
			log.Println("连接mongo数据库成功！",Config.MongoAddress)
		}
	}
	//最大连接池默认为4096
	mgoSession.SetPoolLimit(100)
	return mgoSession.Clone()
}

// mongoDB获取集合对象
func GetMongoCollection(collection string, s func(*mgo.Collection) error ) error{
	session := getSession()
	defer session.Close()
	c := session.DB(Config.MongoDataBaseName).C(collection)
	return s(c)
}

// mongoDB获取collection对象个数
func GetMongoCollectionCount(collectionName string,method func(session *mgo.Collection) (int,error)) (int,error){
	session:= getSession()
	defer session.Close()
	// 数据库名      集合名
	collection := session.DB(Config.MongoDataBaseName).C(collectionName)
	// 使用方法
	return method(collection)
}