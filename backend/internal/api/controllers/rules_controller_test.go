package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/response"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	ctrl := NewRulesController()
	awsRules := r.Group("/aws/rules")
	{
		awsRules.GET("", ctrl.List)
		awsRules.GET("/:service", ctrl.Get)
	}
	return r
}

func TestRulesController_List(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/aws/rules", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rules []response.ServiceRuleResponse
	err := json.Unmarshal(w.Body.Bytes(), &rules)
	assert.NoError(t, err)
	assert.NotEmpty(t, rules)

	// Check for presence of key services
	services := make(map[string]bool)
	for _, rule := range rules {
		services[rule.ServiceName] = true
	}

	assert.True(t, services["VPC"])
	assert.True(t, services["Subnet"])
	assert.True(t, services["EC2"])
	assert.True(t, services["S3"])
	assert.True(t, services["RDS"])
	assert.True(t, services["IAMRole"])
	assert.True(t, services["ECSCluster"])
}

func TestRulesController_Get_VPC(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/aws/rules/VPC", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rule response.ServiceRuleResponse
	err := json.Unmarshal(w.Body.Bytes(), &rule)
	assert.NoError(t, err)

	assert.Equal(t, "VPC", rule.ServiceName)
	// VPC is a top-level resource, valid parents might be empty or specific if we have aggregation
	// Check children - Subnet should be a valid child of VPC (as Subnet has VPC as parent)
	assert.Contains(t, rule.ValidChildren, "Subnet")
	assert.Contains(t, rule.ValidChildren, "InternetGateway")
}

func TestRulesController_Get_Subnet(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/aws/rules/Subnet", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rule response.ServiceRuleResponse
	err := json.Unmarshal(w.Body.Bytes(), &rule)
	assert.NoError(t, err)

	assert.Equal(t, "Subnet", rule.ServiceName)
	assert.Contains(t, rule.ValidParents, "VPC")

	// EC2 should be a valid child of Subnet
	assert.Contains(t, rule.ValidChildren, "EC2")
}

func TestRulesController_Get_NotFound(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/aws/rules/NonExistentService", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
