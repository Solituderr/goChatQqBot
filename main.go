// @title        Golang Service Template
// @version      0.1
// @description  Golang back-end service template, get started with back-end projects quickly
// @BasePath     /api

package main

import (
	"github.com/sirupsen/logrus"
	"go-svc-tpl/app"
	"go-svc-tpl/model"
	"go-svc-tpl/utils"
)

func init() {

}

func main() {
	logrus.SetReportCaller(true)
	model.Init()
	app.InitWebFramework()
	go utils.PushPeripheral()
	go utils.ClearUserInfo()
	app.StartServer()
}
