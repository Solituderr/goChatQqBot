package app

import (
	"go-svc-tpl/app/controller"
)

func addRoutes() {
	e.POST("/", controller.ClassifyReq)
}
