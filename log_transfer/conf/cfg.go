package conf

type LogTransfer struct {
	KafkaCfg `ini:"kafka"`
	ESCfg    `ini:"es"`
}

type KafkaCfg struct {
	Adress string `ini:"address"`
	Topic  string `ini:"topic"`
}

type ESCfg struct {
	Adress   string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
	Nums     int    `ini:"nums"`
}
