package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"log"


	"grpc/common/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc/common/bet"
	"algameserver/controller"
)

// API routes
func API(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/bet", func(c *gin.Context) {
		msg := bet1()
		c.JSON(200, gin.H{
			"message": msg,
		})
	})

	router.GET("/test", func(context *gin.Context) {

		path := fmt.Sprintf("http://127.0.0.1:8881/home")

		tr := &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
		}
		netClient := &http.Client{
			Transport: tr,
			Timeout:   time.Second * 5,
		}

		response, err := netClient.Get(path)
		if err != nil {
			fmt.Printf("err %s", err)
		}

		//程序在使用完回复后必须关闭回复的主体。
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("ReadAll err %s", err)
		}

		context.JSON(200, gin.H{
			"message": string(body),
		})

	})

	router.GET("/db/now", controller.DBnow)
	router.GET("/redis/server", controller.RedisServer)
}

func bet1() string {
	// 連線到遠端 gRPC 伺服器。
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	// 建立新的 BetService 客戶端
	c := bet.NewBetServiceClient(conn)

	r, err := c.Bet(context.Background(), &bet.BetRequest{
		Amount:   12,
		PeriodNo: "123",
		User: &user.User{
			UserID: 22,
			Aid:    44},
	})

	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	return r.Message
}
