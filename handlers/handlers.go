package handlers

import (
	"log"

	"github.com/cerbos/cerbos/client"
	"api/model"
	"github.com/gin-gonic/gin"
)

type configHandler struct {
	host string
}

const (
	username = "cerbos"
	password = "randomHash"
)

func NewConfigHandler(host string) *configHandler {
	return &configHandler{
		host,
	}
}
func (h *configHandler) CheckPolicyHandler(g *gin.Context) {
	payload := &model.CerbosPayload{}
	g.BindJSON(payload)
	cli, err := client.New("dns:///cerbos:3593", client.WithPlaintext())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		g.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}
	CheckPolicy(payload, cli, g)
}

func (h *configHandler) AddPolicyHandler(g *gin.Context) {
	payload := &model.CerbosPayload{}
	g.BindJSON(payload)
	cli, err := client.NewAdminClientWithCredentials(h.host, username, password, client.WithPlaintext())
	if err != nil {
		log.Fatalf("unable to connect to grpc: %v", err)
		g.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	AddPolicy(payload, cli, g)
}
