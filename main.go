package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rawmaterials223/MiniDouyin/controller"
	"github.com/rawmaterials223/MiniDouyin/repository"
	"github.com/rawmaterials223/MiniDouyin/util"
)

func main() {
	go controller.RunMessageServer()

	if err := Init(); err != nil {
		return
	}

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func Init() error {
	if err := repository.DbInit(); err != nil {
		return err
	}
	if err := util.InitLogger(); err != nil {
		return err
	}

	return nil
}
