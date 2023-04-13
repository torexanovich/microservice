package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	cp "gitlab.com/micro/api_gateway/genproto/comment"
	l "gitlab.com/micro/api_gateway/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateComment...
// @Summary CreateComment
// @Description This API for creating a new Comment
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param body body models.CommentRequest true "CommentRequest"
// @Success 200 {object} models.CommentResponse
// @Failure 400 {object} models.Error
// @Router /v1/comment/ [post]
func (h *handlerV1) CreateComment(c *gin.Context) {
	var (
		body        cp.CommentRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	response, err := h.serviceManager.CommentService().WriteComment(context.Background(), &cp.CommentRequest{
		PostId: body.PostId,
		UserId: body.UserId,
		Text:   body.Text,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetComments ...
// @Summary GetComments
// @Description This API for getting Comments detail
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} models.Comments
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/comment/{id} [get]
func (h *handlerV1) GetCommentsForPost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int", l.Error(err))
		return
	}

	response, err := h.serviceManager.CommentService().GetCommentsForPost(context.Background(), &cp.GetAllCommentsRequest{PostId: int64(intId)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get comment by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete Comment
// @Description This API for deleting Comment
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param id path int true "Id"
// @Success 200 {object} models.CommentResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/comment/{id} [delete]
func (h *handlerV1) DeleteComment(c *gin.Context) {
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	newID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.CommentService().DeleteComment(context.Background(), &cp.IdRequest{Id: int64(newID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
