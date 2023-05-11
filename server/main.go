package main

import (
	"github.com/gin-gonic/gin"
	"github.com/proc-moe/aarealbs/server/model"
	"github.com/proc-moe/aarealbs/server/routers"
	"github.com/proc-moe/aarealbs/server/service"
	"github.com/proc-moe/aarealbs/server/utils/cors"
	"github.com/proc-moe/aarealbs/server/utils/klog"
)

// import (
// 	"fmt"
// 	"log"

// 	"github.com/gin-gonic/gin"
// )

// func init() {
// }

// // @title Golang Gin API
// // @version 1.0
// // @description An example of gin
// // @termsOfService https://github.com/EDDYCJY/go-gin-example
// // @license.name MIT
// // @license.url https://github.com/EDDYCJY/go-gin-example/blob/master/LICENSE
// func main() {
// 	route := gin.Default()
// 	route.POST("/api/login", login)
// 	route.Run(":8085")

// }

//	func login(c *gin.Context) {
//		var token Token
//		//var p interface{}
//		if c.BindJSON(&token) != nil {
//			log.Printf("token=%v\n", token.token)
//		}
//		fmt.Printf("%v\n", token)
//	}

type Token struct {
	Token string `form:"token"`
}

func main() {
	klog.Init()
	r := gin.Default()
	r.Use(cors.Cors())
	model.Init()
	routers.Init(r)
	service.MonitorStart()
	r.Run(":9987")
}
