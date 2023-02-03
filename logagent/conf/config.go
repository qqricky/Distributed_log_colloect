package conf

type AppConf struct {
	KafkConf `ini:"kafka"`
	EtcdConf `ini:"etcd"`
}

type KafkConf struct {
	Address     string `ini:"address"`
	Chanmaxsize int    `ini:"chan_max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}

// unused-----------
type TaillogConf struct {
	FileName string `ini:"filename"`
}
