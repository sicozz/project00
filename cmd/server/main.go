package main

import (
	"github.com/sicozz/project00/controller"
)

func main() {
	rc, _ := controller.NewRootController()
	rc.Launch()
}
