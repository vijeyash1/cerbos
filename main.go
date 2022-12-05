package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cerbos/cerbos/client"
)

const (
	username = "cerbos"
	password = "randomHash"
)

// host "dns:///cerbos:3593"
type configHandler struct {
	host string
}

// CheckPolicyHandler checks for the added policy
func (h *configHandler) CheckPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// for creating normal client no username or password
	c, err := client.New("dns:///cerbos:3593", client.WithPlaintext())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		w.Write([]byte(err.Error()))
		return
	}
	//  principle is created
	principal := client.NewPrincipal("someidhere", "admin")
	principal.WithScope("nokia")
	resource := "contact"
	actions := []string{"create", "update", "delete", "read", "list"}

	r1 := client.NewResource(resource, "someidhere")
	// r1.WithScope("nokia")
	//  resource is created
	batch := client.NewResourceBatch()
	batch.Add(r1, actions...)

	// to check resource it needs principle and resouce batch
	//  we are creating the principle and the resource batch above and adding here
	resp, err := c.CheckResources(context.Background(), principal, batch)
	if err != nil {
		log.Fatalf("Failed to check resources: %v", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("response : %v", resp)))

}

// AddPolicyHandler adds a new resource policy
func (h *configHandler) AddPolicyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// for creating admin client username and password required
	c, err := client.NewAdminClientWithCredentials(h.host, username, password, client.WithPlaintext())
	if err != nil {
		log.Fatalf("unable to connect to grpc: %v", err)
		w.Write([]byte(err.Error()))
		return
	}
	// new policy set created
	ps := client.NewPolicySet()

	// actions for user role
	userActions := []string{
		"read",
		"list",
	}

	// actions for admin role
	adminActions := []string{
		"read",
		"list",
		"create",
		"update",
		"delete",
	}

	rr1 := client.NewAllowResourceRule(userActions...).WithRoles("user")

	rr2 := client.NewAllowResourceRule(adminActions...).WithRoles("admin")
	resourcePolicy := client.NewResourcePolicy("contact", "default").AddResourceRules(rr1).AddResourceRules(rr2)
	// here we add resource policy to the policyset
	// resourcePolicy.WithScope("nokia")
	policySet := ps.AddResourcePolicies(resourcePolicy)
	// here using adminclient we are adding the newly created policyset
	err = c.AddOrUpdatePolicy(context.Background(), policySet)
	if err != nil {
		log.Fatalf("unable to add policy: %v", err)

		w.Write([]byte(err.Error()))

	}
	resourcePolicy.WithScope("nokia")
	err = c.AddOrUpdatePolicy(context.Background(), policySet)
	if err != nil {
		log.Fatalf("unable to add policy: %v", err)

		w.Write([]byte(err.Error()))
		return
	}

	log.Println("created the resource policy")
	w.Write([]byte(fmt.Sprintf("------------------------------------Added resource policy -------------------------------------------   \n %+v ", policySet.GetPolicies())))

}

func main() {
	h := configHandler{}
	h.host = "dns:///cerbos:3593"

	http.HandleFunc("/add", h.AddPolicyHandler)
	http.HandleFunc("/check", h.CheckPolicyHandler)

	http.ListenAndServe(":8090", nil)

}
