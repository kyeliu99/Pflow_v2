package flow

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/pflow/backend/internal/flow"
)

type Handlers struct {
	Service flow.Service
}

type createFlowRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Definition  map[string]any    `json:"definition" binding:"required"`
	Metadata    map[string]string `json:"metadata"`
}

type updateFlowRequest struct {
	Description string            `json:"description"`
	Definition  map[string]any    `json:"definition" binding:"required"`
	Metadata    map[string]string `json:"metadata"`
}

func (h Handlers) List(c *gin.Context) {
	flows, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flows)
}

func (h Handlers) Create(c *gin.Context) {
	var req createFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.Service.Create(c.Request.Context(), flow.CreateInput{
		Name:        req.Name,
		Description: req.Description,
		Definition:  req.Definition,
		Metadata:    req.Metadata,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h Handlers) Get(c *gin.Context) {
	result, err := h.Service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if flow.IsNotFound(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handlers) Update(c *gin.Context) {
	var req updateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.Service.Update(c.Request.Context(), flow.UpdateInput{
		ID:          c.Param("id"),
		Description: req.Description,
		Definition:  req.Definition,
		Metadata:    req.Metadata,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if flow.IsNotFound(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}
