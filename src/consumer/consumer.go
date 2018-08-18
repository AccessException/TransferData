package main

import (
	"entity"
	"utils"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"JaegerTracer"
	"time"
	"github.com/Shopify/sarama"
	"encoding/json"
	"logger"
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
)

// 驼峰命名（分区有问题）
var topic = "transfer_data"

type myStruct struct {
	ID 		string 				`json:"_id"`
	entity.Model
}

func main() {
	db ,err := gorm.Open("mysql","demo:demo@tcp(localhost:4000)/demo")
	if err != nil{
		log.Println("***",err)
	}
	// 自动创建表
	db.AutoMigrate(&myStruct{})
	topics := []string{topic}
	var wg = &sync.WaitGroup{}
	wg.Add(5)
	go clusterConsumer(wg, []string{"127.0.0.1:9092"}, topics, "group-1",*db)
	wg.Wait()
}

// 修改默认表名
func (myStruct) TableName() string {
	return utils.Config.TiDBTableName
}

// 支持brokers cluster的消费者
func clusterConsumer(wg *sync.WaitGroup,brokers, topics []string, groupId string,db gorm.DB) {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				mystruct := myStruct{}
				// 修改表名
				mystruct.TableName()
				json.Unmarshal(msg.Value,&mystruct)
				var writeLog map[string]interface{}
				json.Unmarshal(msg.Value,&writeLog)
				logger.WriteLogFile(writeLog,"consumer")
				db.Create(&mystruct)
				tracer, close, _ := JaegerTracer.NewJaegerTracer("Kafka"+mystruct.ID,"127.0.0.1:6831")
				context := make(map[string]string)
				json.Unmarshal(msg.Key,&context)
				spanContext, _ := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(context))
				span := opentracing.StartSpan(
					"liupengKafkaConsumer",
					opentracing.ChildOf(spanContext),
				)
				span.SetTag("mongoId",mystruct.ID)
				span.SetTag("出Kafka时间",time.Now().Format("2006-01-02 15:04:05"))
				span.Finish()
				close.Close()
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}