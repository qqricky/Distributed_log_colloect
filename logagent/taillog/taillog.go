package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"logagent/kafka"
)

//专门从日志文件收集日志并发往kafka

// 一个日志收集的任务
type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// 开启一个新任务
func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init()
	return
}

func (t TailTask) init() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}

	go t.run()
}

func (t TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			//fmt.Printf("tail task%s_%s结束了...\n", t.path, t.topic)
			return
		case line := <-t.instance.Lines: //从tailObj中一行一行读取日志
			//3.2发往kafka
			//先把日志发到一个通道中
			fmt.Printf("get log data from %s success, logdata:%v\n", t.path, line.Text)
			kafka.SendToChan(t.topic, line.Text)
		}
	}
}
