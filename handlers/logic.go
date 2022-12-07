package handlers

import (
	"context"
	"fmt"

	"api/model"

	"github.com/cerbos/cerbos/client"
	"github.com/gin-gonic/gin"
)

func CheckPolicy(datas *model.CerbosPayload, cli client.Client, g *gin.Context) {
	response := model.Response{}
	responses := []model.Response{}
	for _, data := range datas.Policies {
		for i, policy := range data.ResourcePolicy.Rules {
			principal := client.NewPrincipal("idid", policy.Roles)
			// principal.WithScope(data.ResourcePolicy.Scope)
			resource := data.ResourcePolicy.Resource
			actions := policy.Actions
			r1 := client.NewResource(resource, "id")
			if data.ResourcePolicy.Scope != "" {
				r1.WithScope(data.ResourcePolicy.Scope)
			}
			// r1.WithScope(data.ResourcePolicy.Scope)
			batch := client.NewResourceBatch()
			batch.Add(r1, actions...)
			resp, err := cli.CheckResources(context.Background(), principal, batch)
			if err != nil {
				g.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			response = model.Response{
				Response: fmt.Sprintf("response %d : %+v", i, resp),
			}
			responses = append(responses, response)
		}
	}
	g.JSON(200, gin.H{
		"res": responses,
	})
}

func AddPolicy(datas *model.CerbosPayload, cli client.AdminClient, g *gin.Context) {
	response := model.Response{}
	responses := []model.Response{}
	for _, data := range datas.Policies {
		for _, policy := range data.ResourcePolicy.Rules {
			ps := client.PolicySet{}
			actions := policy.Actions
			rr1 := client.NewAllowResourceRule(actions...).WithRoles(policy.Roles)
			resource := data.ResourcePolicy.Resource
			resourcePolicy := client.NewResourcePolicy(resource, data.ResourcePolicy.Version).AddResourceRules(rr1)
			if data.ResourcePolicy.Scope != "" {
				resourcePolicy.WithScope(data.ResourcePolicy.Scope)
			}
			policySet := ps.AddResourcePolicies(resourcePolicy)
			err := cli.AddOrUpdatePolicy(context.Background(), policySet)
			if err != nil {
				g.JSON(400, gin.H{
					"error": err.Error(),
				})
				return

			}
			response = model.Response{
				Response: fmt.Sprintf("response : %+v", policySet.GetPolicies()),
			}
			responses = append(responses, response)
		}
	}
	g.JSON(200, gin.H{
		"Add policy response": responses,
	})
}
