package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
	"time"
)

// 初始化es，接收kafka发来的数据
type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *LogData
)

func Init(address string, chanSize, nums int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return
	}
	fmt.Println("connect to es success")
	ch = make(chan *LogData, chanSize)
	for i := 0; i < nums; i++ {
		go sendToES()
	}

	return
}

func SendToESChan(msg *LogData) {
	ch <- msg
}
func sendToES() {
	for {
		select {
		case msg := <-ch:
			put1, err := client.Index().Index(msg.Topic).Type("xxx").BodyJson(msg).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Indexed student %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}

}
