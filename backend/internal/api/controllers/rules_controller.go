package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/response"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
)

type RulesController struct{}

func NewRulesController() *RulesController {
	return &RulesController{}
}

// List returns all service rules
// @Summary      List all service rules
// @Description  Get all AWS service rules and their constraints
// @Tags         rules
// @Produce      json
// @Success      200  {array}   response.ServiceRuleResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /aws/rules [get]
func (c *RulesController) List(ctx *gin.Context) {
	allRules := getAllRules()
	response := processRules(allRules, "")
	ctx.JSON(http.StatusOK, response)
}

// Get returns rules for a specific service
// @Summary      Get service rules
// @Description  Get rules and constraints for a specific AWS service
// @Tags         rules
// @Produce      json
// @Param        service  path      string  true  "Service Name (e.g. VPC, Subnet)"
// @Success      200      {object}  response.ServiceRuleResponse
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /aws/rules/{service} [get]
func (c *RulesController) Get(ctx *gin.Context) {
	serviceName := ctx.Param("service")
	allRules := getAllRules()

	// We process ALL rules first to correctly build relationships (children/parents),
	// then we filter for the requested service.
	responses := processRules(allRules, serviceName)

	if len(responses) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Return the single service response
	ctx.JSON(http.StatusOK, responses[0])
}

func getAllRules() []rules.ConstraintRecord {
	var allRules []rules.ConstraintRecord
	allRules = append(allRules, rules.DefaultNetworkingRules()...)
	allRules = append(allRules, rules.DefaultComputeRules()...)
	allRules = append(allRules, rules.DefaultStorageRules()...)
	allRules = append(allRules, rules.DefaultDatabaseRules()...)
	allRules = append(allRules, rules.DefaultIAMRules()...)
	allRules = append(allRules, rules.DefaultContainerRules()...)
	return allRules
}

func processRules(constraints []rules.ConstraintRecord, filterService string) []response.ServiceRuleResponse {
	// 1. Group rules by ResourceType
	// 2. Build Parent->Children map inverted from Child->Parent rules

	resourceRules := make(map[string][]rules.ConstraintRecord)
	parentToChildren := make(map[string][]string)

	// First pass: Organize rules and build parent-child relationships
	for _, r := range constraints {
		resourceRules[r.ResourceType] = append(resourceRules[r.ResourceType], r)

		if r.ConstraintType == "requires_parent" || r.ConstraintType == "allowed_parent" {
			parents := strings.Split(r.ConstraintValue, ",")
			for _, p := range parents {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				// r.ResourceType is a child of p
				if !contains(parentToChildren[p], r.ResourceType) {
					parentToChildren[p] = append(parentToChildren[p], r.ResourceType)
				}
			}
		}
	}

	var result []response.ServiceRuleResponse

	// If filterService is provided, we only want to generate response for that service
	// BUT we needed the full pass above to get valid children.

	targetServices := []string{}
	if filterService != "" {
		// keys are case sensitive in map, but search should be case-insensitive
		for k := range resourceRules {
			if strings.EqualFold(k, filterService) {
				targetServices = append(targetServices, k)
			}
		}
	} else {
		for k := range resourceRules {
			targetServices = append(targetServices, k)
		}
	}

	for _, serviceName := range targetServices {
		records := resourceRules[serviceName]
		resp := response.ServiceRuleResponse{
			ServiceName:   serviceName,
			Rules:         []response.RuleDetail{},
			ValidParents:  []string{},
			ValidChildren: []string{},
		}

		// Populate Rules and ValidParents
		for _, r := range records {
			resp.Rules = append(resp.Rules, response.RuleDetail{
				Type:  r.ConstraintType,
				Value: r.ConstraintValue,
			})

			if r.ConstraintType == "requires_parent" || r.ConstraintType == "allowed_parent" {
				parents := strings.Split(r.ConstraintValue, ",")
				for _, p := range parents {
					p = strings.TrimSpace(p)
					if p != "" && !contains(resp.ValidParents, p) {
						resp.ValidParents = append(resp.ValidParents, p)
					}
				}
			}
		}

		// Populate ValidChildren from our pre-calculated map
		if children, ok := parentToChildren[serviceName]; ok {
			resp.ValidChildren = children
		}

		result = append(result, resp)
	}

	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
