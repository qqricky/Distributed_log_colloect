package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var tskMgr *taillogMgr

// tailTask管理者
type taillogMgr struct {
	//logEntry
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry:    logEntryConf, //把当前的日志收集项配置信息保存起来
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}
	for _, logEntey := range logEntryConf {
		//conf: etcd.LogEntry
		//LogEntry.path要收集的文件的路径
		//初始化的时候起了多少个tailTask都要记下来

		tailObj := NewTailTask(logEntey.Path, logEntey.Topic)
		mk := fmt.Sprintf("%s_%s", logEntey.Path, logEntey.Topic)
		tskMgr.tskMap[mk] = tailObj
		//return
	}
	go tskMgr.runwatch()
}

// 监听newConfChan有了新的配置就做处理
func (t *taillogMgr) runwatch() {
	for {
		select {
		case newConf := <-t.newConfChan:
			//释放原有的监听goroutinne，并删除t.tskMap中对应的键值对
			for _, conf := range t.logEntry {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				t.tskMap[mk].cancelFunc()
				delete(t.tskMap, mk)
			}
			//重新开启新传来配置的监听
			for _, logEntey := range newConf {
				tailObj := NewTailTask(logEntey.Path, logEntey.Topic)
				mk := fmt.Sprintf("%s_%s", logEntey.Path, logEntey.Topic)
				tskMgr.tskMap[mk] = tailObj
			}
			t.logEntry = newConf
			fmt.Println("新的配置来了", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 一个函数，向外暴露tskMgr的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
