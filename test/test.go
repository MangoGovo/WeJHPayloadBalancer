package test

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

//func main() {
//	r := gin.Default()
//	r.GET("/api/test", test)
//	_ = r.Run()
//}

type requestData struct {
	MaxDelay    int     `form:"max_delay" bind:"gt=0"`
	FailureRate float64 `form:"failure_rate" bind:"gt=0"`
}

func Test(c *gin.Context) {
	data := requestData{
		MaxDelay:    0,
		FailureRate: 0,
	}
	err := c.ShouldBind(&data)

	if err != nil {
		c.JSON(500, gin.H{"msg": "params error"})
		return
	}

	// 概率失败
	if rand.Float64() < data.FailureRate {
		c.JSON(500, gin.H{"msg": "error"})
		return
	}

	// 随机延时
	delay := 0
	if data.MaxDelay > 0 {
		delay = rand.Intn(data.MaxDelay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	c.JSON(200, gin.H{
		"message": "success",
		"delay":   delay,
	})
}
