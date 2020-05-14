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
	{
		// 返回心跳
		// 格式:
		/*
		   {
		       "status": "alive",
		       "uid": "香港记者号",
		       "queue": {
		           "length": "10", "//": "队列长度",
		           "max": "20",    "//": "最大队列长度"
		       }
		   }
		*/
		Router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "alive",
				"uid":    JudgerID,
			})
		})
	}
	{
		// 提交评测代码以及数据
		// 请求格式:
		/*
		   {
		       "id": "主服务器上的提交ID 此项会在回调时提供",
		       "code": "用户的代码",
		       "limits": {
		           "cpu": "1000", "//": "CPU时间限制 单位为ms. 用户进程将会在相当于2.25倍CPU时间限制的实际时间后被杀死.",
		           "memory": "", "//": "内存使用限制 单位为Byte 实际内存限制会是这里传入的2倍, 请在随后传回的内存用量中自行判断是否MLE",
		       },
		       "data": [
		           "",
		           ""
		       ], "//": "输入数据. 必须为字符串形式."
		   }
		*/
		// 成功提交返回格式:
		/*
		   {
		       "id": "传入的id",
		       "status": "ok"
		   }
		*/
		// 提交失败的返回格式(如评测鸡实例队列已满)
		/*
		   {
		       "id": "传入的id",
		       "status": "queue-full"
		   }
		*/
		Router.POST("/submit", func(c *gin.Context) {
			// TODO
		})
	}
	{
		Router.GET("/test", func(c *gin.Context) {
		})
	}
}
