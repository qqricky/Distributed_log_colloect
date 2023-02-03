package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"logagent/utils"
	"sync"
	"time"
)

//logAgent入口

var (
	cfg = new(conf.AppConf)
)

func main() {
	//0.加载配置文件
	//将配置绑定到cfg结构体
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}

	//1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkConf.Address}, cfg.KafkConf.Chanmaxsize)
	if err != nil {
		fmt.Print("init kafka failed ,err%v\n", err)
		return
	}
	fmt.Println("init kafka seccess")

	//2.初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Print("init etcd failed ,err%v\n", err)
		return
	}
	fmt.Println("init kafka etcd")

	//为了实现每个logagent都拉去自己独有的配置，所以要以自己的IP地址作为区分
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr)
	//2.1 从etcd中获取日志收集项配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed,err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success,%v\n", logEntryConf)
	//2.2 派一个哨兵去监视日志收集项的变化(有变化及时通知我的logAgent实现热加载配置)

	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v\n", index, value)
	}
	//3.收集日志发往kafka
	taillog.Init(logEntryConf)

	newConfChan := taillog.NewConfChan() //从taillog包中获取对外暴露的通道

	var wg sync.WaitGroup
	wg.Add(1)
	//监视key为etcdConfKey的value变化
	go etcd.WatchConf(etcdConfKey, newConfChan) //哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()
}
