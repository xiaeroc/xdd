package models

import (
	"github.com/beego/beego/v2/core/logs"
	"os"
)

var test2 = func(string) {

}

func init() {
	killp()
	for _, arg := range os.Args {
		if arg == "-d" {
			Daemon()
		}
	}
	path, _ := os.Getwd()
	logs.Info("当前%s", ExecPath)
	//path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	ExecPath = path
	logs.Info("当前%s", ExecPath)
	initConfig()
	initDB()
	go initVersion()
	//go initUserAgent()
	initContainer()
	initHandle()
	initFunction()
	//initCron()
	go initTgBot()
	InitReplies()
	initTask()
	//initRepos()
	intiSky()
}
