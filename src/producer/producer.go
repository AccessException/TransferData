package main

import (
	"log"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"gopkg.in/mgo.v2/bson"
	"utils"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/opentracing/opentracing-go"
	"context"
	"JaegerTracer"
	"entity"
	"encoding/json"
	"strconv"
	"logger"
	"time"
)

var topic = "transferData"

type MyStruct struct {
	ID 		bson.ObjectId 		`json:"_id"`
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
	collections,err := utils.FindAll(bson.M{},utils.Config.MongoCollectionName)
	if err != nil{
		log.Println("获取数据错误")
	}
	for i:=0;i<len(collections);i++ {
		var myStruct MyStruct
		var model entity.Model
		content, err := json.Marshal(collections[i])
		if err != nil {
			log.Println("转字节错误", err)
		}
		json.Unmarshal(content,&struct {
			*MyStruct
			*entity.Model
			}{&myStruct, &model})

		body,err := json.Marshal(struct {
			entity.Model
			ID string `json:"_id"`
		}{
			Model: model,
			ID: myStruct.ID.Hex(),
		})
		var writeLog map[string]interface{}
		json.Unmarshal(body,&writeLog)
		logger.WriteLogFile(writeLog,"producer")
		header := make(map[string]string)
		tracer, close, _ := JaegerTracer.NewJaegerTracer("Kafka"+myStruct.ID.Hex(), "127.0.0.1:6831")
		span := tracer.StartSpan("liupengKafkaProducer")
		md := make(map[string]string)
		err = tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
		span.SetTag("策略","10%取样")
		span.SetTag("mongoId",myStruct.ID.Hex())
		span.SetTag("进Kafka时间",time.Now().Format("2006-01-02 15:04:05"))
		header = md
		span.Finish()
		close.Close()
		msg := &broker.Message{
			Header: header,
			Body:   body,
		}
		if err := myBroker.Publish(topic, msg); err == nil {
			log.Println("发送成功！"+strconv.Itoa(i))
		} else {
			log.Println("发送失败！", err.Error())
		}

	}
}