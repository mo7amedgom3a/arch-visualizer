package response

type ServiceRuleResponse struct {
	ServiceName   string       `json:"service_name"`
	Description   string       `json:"description,omitempty"`
	Rules         []RuleDetail `json:"rules"`
	ValidParents  []string     `json:"valid_parents,omitempty"`
	ValidChildren []string     `json:"valid_children,omitempty"`
}

type RuleDetail struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}
