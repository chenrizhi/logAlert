package main

import (
	"encoding/json"
	"github.com/astaxie/beego/config/yaml"
	logs "github.com/sirupsen/logrus"
)

var (
	appConfig map[string]map[string]string
)

func initConfig(filename string) (err error) {
	conf, err := yaml.ReadYmlReader(filename)
	if err != nil {
		logs.Error("read config file failed, err:", err)
	}
	confStr, _ := json.Marshal(conf)
	err = json.Unmarshal(confStr, &appConfig)
	if err != nil {
		logs.Error("init config failed, err:", err)
	}
	logs.Info("init config success")
	return
}