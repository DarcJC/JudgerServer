package router

import (
	"JudgerServer/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Router 路由
var Router *gin.Engine

// JudgerID 判题器的ID
var JudgerID string

func init() {
	v, err := config.GetConfig("UNIQUE_ID")
	if err != nil {
		panic(err)
	}
	if v == "" {
		JudgerID = uuid.New().String()
	} else {
		JudgerID = v
	}

	Router = gin.Default()
	Router.Use(BasicAuth())
	{
		// 返回心跳
		// 格式:
		/*
		   {
		       "status": "alive",
		       "uid": "香港记者号"
		   }
		*/
		Router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "alive",
				"uid":    JudgerID,
			})
		})
	}
}
