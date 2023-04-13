package v1

import (
	"context"
	"net/http"
	"strconv"

	pp "gitlab.com/micro/api_gateway/genproto/post"
	l "gitlab.com/micro/api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreatePost...
// @Summary CreatePost
// @Description This API for creating a new post
// @Tags Post
// @Accept  json
// @Produce  json
// @Param body body models.PostRequest true "postRequest"
// @Success 200 {object} models.GetPostResponse
// @Failure 400 {object} models.Error
// @Router /v1/post [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        pp.PostRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	response, err := h.serviceManager.PostService().CreatePost(context.Background(), &pp.PostRequest{
		Title:       body.Title,
		Description: body.Description,
		UserId:      body.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPost ...
// @Summary GetPost
// @Description This API for getting Post detail
// @Tags Post
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} models.GetPostResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/post/{id} [get]
func (h *handlerV1) GetPost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int", l.Error(err))
		return
	}

	response, err := h.serviceManager.PostService().GetPostById(context.Background(), &pp.IdRequest{Id: int64(intID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPostByUserId ...
// @Summary GetPostByUserId
// @Description This API for getting posts by user id
// @Tags Post
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} models.Posts
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/posts/{id} [get]
func (h *handlerV1) GetPostByUserId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	response, err := h.serviceManager.PostService().GetPostByUserId(context.Background(), &pp.IdUser{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get posts by user id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update Post
// @Description This API for updating Post
// @Tags Post
// @Accept  json
// @Produce  json
// @Param body body models.UpdatePostReq true "UpdatePostReq"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/post/{id} [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        pp.UpdatePostRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	response, err := h.serviceManager.PostService().UpdatePost(context.Background(), &pp.UpdatePostRequest{
		Title:       body.Title,
		Description: body.Description,
		Id:          body.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete Post
// @Description This API for deleting Post
// @Tags Post
// @Accept  json
// @Produce  json
// @Param id path int true "Id"
// @Success 200 {object} models.GetPostResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/post/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to atoi id: ", l.Error(err))
	}

	response, err := h.serviceManager.PostService().DeletePost(context.Background(), &pp.IdRequest{Id: int64(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to delete company: ", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)

}
