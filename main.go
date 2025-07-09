// @title        Golang Service Template
// @version      0.1
// @description  Golang back-end service template, get started with back-end projects quickly
// @BasePath     /api

package main

import (
	"go-svc-tpl/app"
	"go-svc-tpl/cronjob"
	"go-svc-tpl/model"
)

func init() {

}

func main() {
	model.Init()
	app.InitLagrangeBot()
	go cronjob.PushPeripheral()
	go cronjob.ClearUserInfo()
	app.StartLagrangeBot()
}
