package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/micro/api_gateway/api/handlers/models"
)

// @Summary Add Policy User
// @Description Add policy for user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Succes 200 {object} user.Empty
// @Router /v1/admin/add/policy [post]
func (h *handlerV1) AddRoleUser(c *gin.Context) {
	body := models.Policy{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	ok, err := h.casbin.AddPolicy(body.User, body.Domain, body.Action)
	if err != nil {
		fmt.Println(">>>", err)
	}

	h.casbin.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, models.Empty{})
}

// @Summary Remove Policy User
// @Description Remove policy for user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Succes 200 {object} user.Empty
// @Router /v1/admin/remove/policy [post]
func (h *handlerV1) RemoveRoleUser(c *gin.Context) {
	body := models.Policy{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	ok, err := h.casbin.RemovePolicy(body.User, body.Domain, body.Action)
	if err != nil {
		fmt.Println(">>>", err)
	}

	h.casbin.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, models.Empty{})
}

// @Summary Get Policy User
// @Description Get policy for user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Router /v1/admin/get/policy [get]
func (h *handlerV1) GetPolicy(c *gin.Context) {
	data := h.casbin.GetPolicy()
	fmt.Println(data)
	c.JSON(http.StatusOK, data)
}


// @Summary Update Policy User
// @Description Update policy for user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.UpdatePol true "Policy"
// @Succes 200 {object} user.Empty
// @Router /v1/admin/update/policy [put]
func (h *handlerV1) UpdatePolicy(c *gin.Context) {
	body := models.UpdatePol{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	new := []string{body.NewPolicy.User, body.NewPolicy.Domain, body.NewPolicy.Action}
	old := []string{body.OldPolicy.User, body.OldPolicy.Domain, body.OldPolicy.Action}
	ok, err := h.casbin.UpdatePolicy(old, new)
	if err != nil {
		fmt.Println(err)
	}
	h.casbin.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, models.Empty{})
}