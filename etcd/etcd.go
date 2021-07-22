package etcd

import (
	"context"
	"encoding/json"
	logs "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

var (
	client *clientv3.Client
)

type LogAlert struct {
	Type     string        `json:"type"`     // 告警类型
	Query    string        `json:"query"`    // 告警条件
	Interval time.Duration `json:"interval"` // 查询间隔
	NumEvent int           `json:"numEvent"` // 触发告警条数
}

func InitEtcd(addrs string, timeout time.Duration) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(addrs, ","),
		DialTimeout: timeout})
	if err != nil {
		return
	}
	logs.Info("init etcd success")
	return
}

func GetConf(key string) (logAlert []*LogAlert, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}
	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, &logAlert)
		if err != nil {
			return
		}
	}
	return
}

func WatchConf(key string, newConfCh chan<- []*LogAlert) {
	ch := client.Watch(context.Background(), key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			logs.Info("Type:%v key:%v value:%v", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			var newConf []*LogAlert
			if evt.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logs.Warn("unmarshal failed, err:%v", err)
					continue
				}
			}
			logs.Info("get new conf:%v", newConf)
			newConfCh <- newConf
		}
	}
}
