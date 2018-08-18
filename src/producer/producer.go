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
	"github.com/Shopify/sarama"
	"strconv"
)

var topic = "transfer_data"

type MyStruct struct {
	ID 		bson.ObjectId 		`json:"_id"`
}

// Kafka 分区在 /usr/local/etc/kafka/server.properties中设置
// num.partitions=5 创建5个分区 修改后重新启功Zookeeper Kafka
func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	producer, _ := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	count ,_:= utils.Count(nil,utils.Config.MongoCollectionName)
	// 修改成分页查询（有待修改分页方法）
	page := 0
	if count%1000 == 0 {
		page = count/1000
	}else{
		page = count/1000+1
	}
	 for j:=0;j<page;j++{
		collections,err := utils.PagingFind(j*1000,1000,nil,utils.Config.MongoCollectionName)
		//collections,err := utils.FindAll(nil,utils.Config.MongoCollectionName)
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

			value,err := json.Marshal(struct {
				entity.Model
				ID string `json:"_id"`
			}{
				Model: model,
				ID: myStruct.ID.Hex(),
			})
			var writeLog map[string]interface{}
			json.Unmarshal(value,&writeLog)
			logger.WriteLogFile(writeLog,"producer")
			tracer, close, _ := JaegerTracer.NewJaegerTracer("Kafka"+myStruct.ID.Hex(), "127.0.0.1:6831")
			span := tracer.StartSpan("liupengKafkaProducer")
			context := make(map[string]string)
			err = tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(context))
			span.SetTag("策略","10%取样")
			span.SetTag("mongoId",myStruct.ID.Hex())
			span.SetTag("进Kafka时间",time.Now().Format("2006-01-02 15:04:05"))
			span.Finish()
			close.Close()
			key,_:=json.Marshal(context)
			msg.Key = sarama.ByteEncoder(key)
			msg.Value = sarama.ByteEncoder(value)
			partition,offset,err := producer.SendMessage(msg)
			if err == nil{
				log.Println("发送成功！"+strconv.Itoa(i),"kafka分区为：",partition,"偏移量为：",offset,"message：",myStruct.ID.Hex())
			}else{
				log.Println("发送失败！")
			}
		}
	}

	producer.Close()
}