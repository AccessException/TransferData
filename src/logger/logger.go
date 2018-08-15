package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/lestrrat-go/file-rotatelogs"
	"time"
	"encoding/json"
)

var Log *logrus.Logger

// 日志文件分隔
func WriteLogFile(content map[string]interface{},fileName string){
	fields := logrus.Fields{}
	logFields := content
	logFields["@timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	fields = logFields
	logf, err := rotatelogs.New(
		// 修改路径
		"/Users/apple/go/TransferData/src/logger/"+fileName+".%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	logrus.SetFormatter(&LogFormat{})
	logrus.SetOutput(logf)
	logrus.WithFields(fields).Info()
	if err != nil{
		logrus.Println(err.Error())
	}
}

type LogFormat struct {}

func (f *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		logrus.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil, err
	}
	return append(serialized, '\n'), nil
}
