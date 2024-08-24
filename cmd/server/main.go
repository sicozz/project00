package main

import (
	"github.com/sicozz/project00/controller"
	"github.com/sicozz/project00/utils"
)

func main() {
	utils.InitLog(utils.DEFAULT_LOG_FILE)
	rc, _ := controller.NewRootController()
	rc.Launch()
}
