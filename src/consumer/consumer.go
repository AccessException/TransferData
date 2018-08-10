package main

import (
	"fmt"
	"log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"gopkg.in/mgo.v2/bson"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"entity"
	"utils"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-plugins/broker/kafka"
	"context"
	"github.com/opentracing/opentracing-go"
	"JaegerTracer"
)

var (
	topic = "loadData"
)

type myStruct struct {
	entity.Model
}

func main() {
	var brokerURLs = []string{"127.0.0.1:9092"}
	myBroker:=kafka.NewBroker(func(o *broker.Options) {
		o.Addrs = brokerURLs
		o.Context = context.Background()
	})
	cmd.Init()
	if err := myBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := myBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	db ,err := gorm.Open("mysql","demo:demo@tcp(localhost:4000)/demo")
	//db ,err := gorm.Open("mysql","root:root@tcp(localhost:3306)/book")
	if err != nil{
		log.Println("***",err)
	}
	// 自动创建表
	db.AutoMigrate(&myStruct{})
	mystruct := myStruct{}
	// 修改表名
	mystruct.TableName()
	_ , err = myBroker.Subscribe(topic, func(p broker.Publication) error {
		tracer, _, _ := JaegerTracer.NewJaegerTracer("KafkaConsumer","127.0.0.1:6831")
		spanContext, _ := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(p.Message().Header))
		span := tracer.StartSpan(
			"liupengKafkaConsumer",
			opentracing.ChildOf(spanContext),
		)
		span.Finish()
		var collection interface{}
		bson.Unmarshal(p.Message().Body,&collection)
		m := entity.Model{}
		mutable := reflect.ValueOf(&m).Elem()
		for i,v:= range collection.(bson.M){
			if i != "_id"{
				mutable.FieldByName(strings.ToUpper(i)).Set(reflect.ValueOf(v))
			}
		}
		s := myStruct{m}
		db.Create(&s)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	select {}
}

// 修改默认表名
func (myStruct) TableName() string {
	return utils.Config.TiDBTableName
}