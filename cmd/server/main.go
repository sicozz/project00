package main

import (
	"fmt"

	"github.com/sicozz/project00/controller"
	"github.com/sicozz/project00/utils"
)

func main() {
	rc, err := controller.NewRootController()
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to start app: %v", err))
		return
	}
	rc.Launch()
}
