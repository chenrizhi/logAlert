package main

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"logAlert/etcd"
	"os"
	"time"
)

var (
	conf string

	rootCmd = &cobra.Command{
		Use: "logalert",
		Run: func(cmd *cobra.Command, args []string) {
			if len(conf) == 0 {
				_ = cmd.Help()
				return
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&conf, "conf", "c", "", "config")
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	execute()

	logs.Info("start logalert")

	err := initConfig(conf)
	if err != nil {
		panic(err)
	}
	err = initLog()
	if err != nil {
		logs.Error("init log failed, err", err)
	}
	err = etcd.InitEtcd(appConfig["etcd"]["address"], time.Second*2)
	if err != nil {
		logs.Error("init etcd failed, err:", err)
		panic(err)
	}
	err = run()
	if err != nil {
		logs.Error("run server failed, err:", err)
	}
}
