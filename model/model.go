package model

type CerbosPayload struct {
	PolicyKind string `json:"policyKind"`
	Policies   []struct {
		APIVersion     string `json:"apiVersion"`
		ResourcePolicy struct {
			Version  string `json:"version"`
			Resource string `json:"resource"`
			Scope    string `json:"scope"`
			Rules    []struct {
				Roles   string   `json:"roles"`
				Actions []string `json:"actions"`
				Effect  string   `json:"effect"`
			} `json:"rules"`
		} `json:"resourcePolicy"`
	} `json:"policies"`
}

type Response struct {
	Response string
}
