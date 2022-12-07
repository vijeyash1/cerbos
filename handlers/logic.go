package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/cerbos/cerbos/client"
	"api/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckPolicy(datas *model.CerbosPayload, cli client.Client, g *gin.Context) {
	response := model.Response{}
	responses := []model.Response{}
	for _, data := range datas.Policies {
		for _, policy := range data.ResourcePolicy.Rules {
			principal := client.NewPrincipal(uuid.NewString(), policy.Roles)
			// principal.WithScope(data.ResourcePolicy.Scope)
			resource := data.ResourcePolicy.Resource
			actions := policy.Actions
			r1 := client.NewResource(resource, uuid.NewString())
			batch := client.NewResourceBatch()
			batch.Add(r1, actions...)
			resp, err := cli.CheckResources(context.Background(), principal, batch)
			if err != nil {
				log.Fatalf("Failed to check resources: %v", err)
				response = model.Response{
					Response: "",
					Errors:   err.Error(),
				}
				responses = append(responses, response)
			}
			response = model.Response{
				Response: fmt.Sprintf("response : %v", resp),
				Errors:   "",
			}
			responses = append(responses, response)
		}
	}
	g.JSON(200, gin.H{
		"Check Responses": responses,
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
			resourcePolicy := client.NewResourcePolicy(resource, "default").AddResourceRules(rr1)
			// resourcePolicy.WithScope(data.ResourcePolicy.Scope)
			policySet := ps.AddResourcePolicies(resourcePolicy)
			err := cli.AddOrUpdatePolicy(context.Background(), policySet)
			if err != nil {
				response = model.Response{
					Response: "",
					Errors:   err.Error(),
				}
				responses = append(responses, response)

			}

			response = model.Response{
				Response: fmt.Sprintf("response : %+v", policySet.GetPolicies()),
				Errors:   "",
			}
			responses = append(responses, response)
		}
	}
	g.JSON(200, gin.H{
		"Check Responses": responses,
	})
}
