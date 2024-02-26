package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/services"
	"net/http"
	"strconv"
)

type Client struct{}
type ClientsHandler struct {
	service *services.Service
}

func NewClientHandler(service *services.Service) *ClientsHandler {
	return &ClientsHandler{service: service}
}

func (c *ClientsHandler) CreateTransaction(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if id > 5 || id < 1 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	t := &data.Transaction{}
	if err := ctx.ShouldBindJSON(t); err != nil {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if !t.Validate() {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	responseTransaction, err := c.service.CreateTransaction(ctx, t, id)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientFunds) {
			ctx.Status(http.StatusUnprocessableEntity)
			return
		}
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx.IndentedJSON(http.StatusOK, responseTransaction)
}

func (c *ClientsHandler) CreateStatement(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if id > 5 || id < 1 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	client := c.service.CreateStatement(ctx, id)
	ctx.IndentedJSON(http.StatusOK, client)
}
