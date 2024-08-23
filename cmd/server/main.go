package main

import (
	rootcontroller "github.com/sicozz/project00/root_controller"
	"github.com/sicozz/project00/utils"
)

func main() {
	utils.InitLog(utils.DEFAULT_LOG_FILE)
	rc, _ := rootcontroller.NewRootController()
	rc.Launch()
}
