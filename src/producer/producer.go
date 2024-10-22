package main

import (
	"log"
	"gopkg.in/mgo.v2/bson"
	"utils"
	"github.com/opentracing/opentracing-go"
	"JaegerTracer"
	"time"
	"entity"
	"encoding/json"
	"logger"
	"strconv"
	//"github.com/Shopify/sarama"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/micro/go-micro/broker"
	"context"
	"github.com/micro/go-micro/cmd"
	"github.com/Shopify/sarama"
)

var topic = "transfer_data"

type MyStruct struct {
	ID 		bson.ObjectId 		`json:"_id"`
}


// 使用 micro broker
func main() {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
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
	count ,_:= utils.Count(nil,utils.Config.MongoCollectionName)
	// 修改成分页查询（有待修改分页方法）
	page := 0
	if count%1000 == 0 {
		page = count/1000
	}else{
		page = count/1000+1
	}
	for j:=0;j<page;j++ {
		collections, err := utils.PagingFind(j*1000, 1000, nil, utils.Config.MongoCollectionName)
		//collections,err := utils.FindAll(nil,utils.Config.MongoCollectionName)
		if err != nil {
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
}



// Kafka 分区在 /usr/local/etc/kafka/server.properties中设置
// num.partitions=5 创建5个分区 修改后重新启功Zookeeper Kafka
// 通过命令修改分区 kafka-topics --zookeeper 127.0.0.1:2181  --alter --topic transfer_data --partitions 4
// 消费者启动会创建topic，partition
//func main() {
//	//设置配置
//	config := sarama.NewConfig()
//	//等待服务器所有副本都保存成功后的响应
//	config.Producer.RequiredAcks = sarama.WaitForAll
//	////随机的分区类型
//	config.Producer.Partitioner = sarama.NewRandomPartitioner
//	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
//	config.Producer.Return.Successes = true
//	config.Version = sarama.V1_0_0_0
//	msg := &sarama.ProducerMessage{}
//	msg.Topic = topic
//	producer, _ := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
//	count ,_:= utils.Count(nil,utils.Config.MongoCollectionName)
//	// 修改成分页查询（有待修改分页方法）
//	page := 0
//	if count%1000 == 0 {
//		page = count/1000
//	}else{
//		page = count/1000+1
//	}
//	for j:=0;j<page;j++{
//		collections,err := utils.PagingFind(j*1000,1000,nil,utils.Config.MongoCollectionName)
//		//collections,err := utils.FindAll(nil,utils.Config.MongoCollectionName)
//		if err != nil{
//			log.Println("获取数据错误")
//		}
//		for i:=0;i<len(collections);i++ {
//			var myStruct MyStruct
//			var model entity.Model
//			content, err := json.Marshal(collections[i])
//			if err != nil {
//				log.Println("转字节错误", err)
//			}
//			json.Unmarshal(content,&struct {
//				*MyStruct
//				*entity.Model
//			}{&myStruct, &model})
//
//			value,err := json.Marshal(struct {
//				entity.Model
//				ID string `json:"_id"`
//			}{
//				Model: model,
//				ID: myStruct.ID.Hex(),
//			})
//			var writeLog map[string]interface{}
//			json.Unmarshal(value,&writeLog)
//			logger.WriteLogFile(writeLog,"producer")
//			tracer, close, _ := JaegerTracer.NewJaegerTracer("Kafka"+myStruct.ID.Hex(), "127.0.0.1:6831")
//			span := tracer.StartSpan("liupengKafkaProducer")
//			context := make(map[string]string)
//			err = tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(context))
//			span.SetTag("策略","10%取样")
//			span.SetTag("mongoId",myStruct.ID.Hex())
//			span.SetTag("进Kafka时间",time.Now().Format("2006-01-02 15:04:05"))
//			span.Finish()
//			close.Close()
//			key,_:=json.Marshal(context)
//			msg.Key = sarama.ByteEncoder(key)
//			msg.Value = sarama.ByteEncoder(value)
//			partition,offset,err := producer.SendMessage(msg)
//			if err == nil{
//				log.Println("发送成功！"+strconv.Itoa(i),"kafka分区为：",partition,"偏移量为：",offset,"message：",myStruct.ID.Hex())
//			}else{
//				log.Println("发送失败！")
//			}
//		}
//	}
//
//	producer.Close()
//}