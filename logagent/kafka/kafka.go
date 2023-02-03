package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type LogData struct {
	topic string
	data  string
}

// 往kafka写日志的模块
var (
	client      sarama.SyncProducer //声明一个全局的的连接kafka的生产者client
	logDataChan chan *LogData
)

func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	logDataChan = make(chan *LogData, maxSize)
	//开启后台的goroutine从通道中取数据发往kafka
	go sendTokafka()
	return
}

// 把日志发送到一个内部的chan中
func SendToChan(topic, data string) {
	msg := &LogData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// 真正将日志信息发送到kafka
func sendTokafka() {
	for {
		select {
		case ld := <-logDataChan:
			//构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			//发送到kafka
			pid, offset, err := client.SendMessage(msg)
			fmt.Println("send data to kafka:", ld.data)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
