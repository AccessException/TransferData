package main

import (
	"log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"gopkg.in/mgo.v2/bson"
	"utils"
	"time"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/opentracing/opentracing-go"
	"context"
	"JaegerTracer"
)

var (
	topic = "loadData"
)

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

	collection,err := utils.FindById(bson.ObjectIdHex("5a715f731fe15a1d0b0a0f3c"),utils.Config.MongoCollectionName)
	if err != nil{
		log.Println("获取数据错误")
	}
	body,err := bson.Marshal(collection)
	if err != nil{
		log.Println("转字节错误",err)
	}
	header := make(map[string]string)
	tracer, _, _ := JaegerTracer.NewJaegerTracer("KafkaProducer","127.0.0.1:6831")
	liuSpan := tracer.StartSpan("liupengKafkaProducer")
	md := make(map[string]string)
	err = tracer.Inject(liuSpan.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	header = md
	liuSpan.Finish()
	msg:=&broker.Message{
		Header:header,
		Body:body,
	}
	tick := time.NewTicker(time.Second)
	for _ = range tick.C {
		if err := myBroker.Publish(topic,msg); err==nil{
			log.Println("发送成功！")
			break
		}else{
			log.Println("发送失败！",err.Error())
		}
	}
}