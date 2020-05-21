package router

import (
	"JudgerServer/config"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var secret string

func init() {
	s, err := config.GetConfig("SECRET")
	if err != nil {
		panic(err)
	}
	secret = s
	if secret == "" {
		secret = RandString(64)
	}
}

// RandString 产生随机字符串
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// BasicAuth 基本鉴权
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"errmsg": "need-authorization",
			})
			return
		}

		if token != secret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errmsg": "authorization-failed",
			})
			return
		}

		c.Next()
	}
}
