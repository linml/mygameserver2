package main

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"game_lib/logging"
	"game_lib/dbconn"
	_ "algameserver/core/services"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"algameserver/core/wsrouter"
	"algameserver/routes"
	"github.com/gorilla/websocket"
	)

var (
	version    string
	date       string
	commit     string

)

func readconfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logging.L().Panic(err.Error())
	}
}

func connectDB() {
	key := "maria.read"
	if !viper.IsSet(key) {
		logging.L().Fatal("config key not found", zap.String("key", key))
	}
	dbconn.OpenRead(viper.GetString(key))

	key = "maria.write"
	if !viper.IsSet(key) {
		logging.L().Fatal("config key not found", zap.String("key", key))
	}
	dbconn.OpenWrite(viper.GetString(key))
}

func connectRedis() {
	key := "redis.addr"
	if !viper.IsSet(key) {
		logging.L().Fatal("config key not found", zap.String("key", key))
	}

	dbkey := "redis.db"
	if !viper.IsSet(key) {
		logging.L().Fatal("config key not found", zap.String("key", dbkey))
	}

	dbconn.RedisDial(viper.GetString(key), viper.GetInt(dbkey))
}

func main() {
	logging.L().Info("info",
		zap.String("version", version),
		zap.String("commit:", commit),
		zap.String("date:", date),
	)

	readconfig()
	connectDB()
	connectRedis()



	logging.L().Info("run service")

	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
		wsrouter.Debug = true
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	engine.LoadHTMLFiles("index.html","index1.html")

	engine.GET("/home", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})



	engine.GET("/home1", func(c *gin.Context) {
		c.HTML(200, "index1.html", nil)
	})

	// websocket router
	wsRouter := wsrouter.New()
	routes.WebsocketRoutes(wsRouter)

	// sets routes
	routes.API(engine)
	routes.WebsocketURL(engine, wsRouter)

	port := viper.GetString("service.port")

	//go wsRouter.BroadcastAction()

	engine.Run(fmt.Sprintf(":%s", port))
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}
