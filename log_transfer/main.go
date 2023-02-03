package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_transfer/conf"
	"log_transfer/es"
	"log_transfer/kafka"
)

//将日志数据从kafka取出来发往es

func main() {
	//0.加载配置文件
	var cfg = new(conf.LogTransfer)

	err := ini.MapTo(cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Println("init conf failed, err:\n", err)
		return
	}
	fmt.Printf("Cfg%v\n", cfg)
	//初始化es
	err = es.Init(cfg.ESCfg.Adress, cfg.ESCfg.ChanSize, cfg.ESCfg.Nums)
	if err != nil {
		fmt.Printf("init es client failed,err:%v\n", err)
		return
	}
	fmt.Println("init es success.\n")
	//1.初始化kafka
	//1.1创建分区的消费者
	//1.2每个分区的消费者分别取出数据，通过SendToES()将数据发往ES
	kafka.Init([]string{cfg.KafkaCfg.Adress}, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed,err:%v\n", err)
		return
	}
	select {}
}
