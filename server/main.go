package main

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
import (
	"fmt"
	"time"
)

type Token struct {
	Token string `form:"token"`
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("nioha")
			}
		}
	}()
	// klog.Init()
	// r := gin.Default()
	// r.Use(cors.Cors())
	// model.Init()
	// routers.Init(r)
	// go service.MonitorStart()
	// r.Run(":9987")
}
