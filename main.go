package main

import (
	"LoadBalance/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func TestRequest(c *gin.Context) {
	threads := []string{
		"http://localhost:8080/api/test?failure_rate=0.5&max_delay=1000",
		"http://localhost:8080/api/test?failure_rate=0.3&max_delay=500",
		"http://localhost:8080/api/test?failure_rate=0.2&max_delay=100",
		"http://localhost:8080/api/test?failure_rate=0.01&max_delay=20",
	}
	data, statusCode := test.Get(threads[0])
	ret, _ := json.Marshal(data)
	c.JSON(statusCode, ret)
}
func main() {
	// 模拟接口
	r4Test := gin.Default()
	r4Test.GET("api/test", test.Test)
	_ = r4Test.Run(":8080")

	// 正式接口
	r := gin.Default()
	r.GET("/api/ping", TestRequest)
	_ = r.Run(":8081")

}
