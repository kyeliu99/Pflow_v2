package workorder

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kyeliu99/Pflow_v2/backend/internal/workorder"
)

type Handlers struct {
	Service workorder.Service
}

type createRequest struct {
	FlowID   string            `json:"flowId" binding:"required"`
	Title    string            `json:"title" binding:"required"`
	Assignee string            `json:"assignee"`
	Payload  map[string]any    `json:"payload"`
	Metadata map[string]string `json:"metadata"`
}

func (h Handlers) List(c *gin.Context) {
	items, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h Handlers) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.Service.Create(c.Request.Context(), workorder.CreateInput{
		FlowID:   req.FlowID,
		Title:    req.Title,
		Assignee: req.Assignee,
		Payload:  req.Payload,
		Metadata: req.Metadata,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h Handlers) Get(c *gin.Context) {
	item, err := h.Service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if workorder.IsNotFound(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h Handlers) Retry(c *gin.Context) {
	if err := h.Service.Retry(c.Request.Context(), c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if workorder.IsNotFound(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
