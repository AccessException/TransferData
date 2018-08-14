package main

import (
	"fmt"
	"log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/jinzhu/gorm"
	"entity"
	"utils"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-plugins/broker/kafka"
	"context"
	"github.com/opentracing/opentracing-go"
	"JaegerTracer"
	"encoding/json"
	"time"
	"github.com/hashicorp/consul/logger"
)

var topic = "transferData"


type myStruct struct {
	ID 		string 				`json:"_id"`
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
	if err != nil{
		log.Println("***",err)
	}
	// 自动创建表
	db.AutoMigrate(&myStruct{})
	mystruct := myStruct{}
	// 修改表名
	mystruct.TableName()
	_ , err = myBroker.Subscribe(topic, func(p broker.Publication) error {
		json.Unmarshal(p.Message().Body,&mystruct)
		db.Create(&mystruct)
		tracer, _, _ := JaegerTracer.NewJaegerTracer("Kafka"+mystruct.ID,"127.0.0.1:6831")
		spanContext, _ := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(p.Message().Header))
		span := opentracing.StartSpan(
			"liupengKafkaConsumer",
			opentracing.ChildOf(spanContext),
		)
		span.SetTag("mongoId",mystruct.ID)
		span.SetTag("出Kafka时间",time.Now().Format("2006-01-02 15:04:05"))
		span.Finish()
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