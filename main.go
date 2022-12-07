package main

import (
	"api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	h := handlers.NewConfigHandler("dns:///cerbos:3593")
	r := gin.Default()
	r.POST("/add", h.AddPolicyHandler)
	r.POST("/check", h.CheckPolicyHandler)
	r.Run()
}
