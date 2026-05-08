package handlers

import (
	"net/http"

	subscrusecase "github.com/IvanKuchsh-600/subtracker/internal/usecase/subscription"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	usecase subscrusecase.Usecase
}

func NewSubscriptionHandler(usecase subscrusecase.Usecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: usecase}
}

// // CreateSubscription godoc
// // @Summary Create new subscription
// // @Tags subscriptions
// // @Accept json
// // @Produce json
// // @Param request body models.CreateSubscriptionRequest true "Subscription data"
// // @Success 201 {object} models.Subscription
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req subscriptionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.usecase.Create(c, subscrusecase.CreateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newSubscriptionDTO(created))
}

// // GetSubscription godoc
// // @Summary Get subscription by ID
// // @Tags subscriptions
// // @Produce json
// // @Param id path string true "Subscription ID"
// // @Success 200 {object} models.Subscription
// // @Failure 404 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /subscriptions/{id} [get]
// func (h *SubscriptionHandler) GetByID(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
// 		return
// 	}

// 	sub, err := h.service.GetByID(c.Request.Context(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if sub == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, sub)
// }

// // UpdateSubscription godoc
// // @Summary Update subscription
// // @Tags subscriptions
// // @Accept json
// // @Produce json
// // @Param id path string true "Subscription ID"
// // @Param request body models.UpdateSubscriptionRequest true "Update data"
// // @Success 200 {object} map[string]string
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req subscriptionMutationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.usecase.Update(c, id, subscrusecase.UpdateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newSubscriptionDTO(updated))
}

// // DeleteSubscription godoc
// // @Summary Delete subscription
// // @Tags subscriptions
// // @Produce json
// // @Param id path string true "Subscription ID"
// // @Success 200 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /subscriptions/{id} [delete]
// func (h *SubscriptionHandler) Delete(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
// 		return
// 	}

// 	if err := h.service.Delete(c.Request.Context(), id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "subscription deleted successfully"})
// }

// ListSubscriptions godoc
// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} models.SubscriptionsListResponse
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
// func (h *SubscriptionHandler) List(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

// 	resp, err := h.service.List(c.Request.Context(), page, pageSize)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

// GetTotalCost godoc
// @Summary Calculate total cost of active subscriptions for period
// @Tags subscriptions
// @Produce json
// @Param from_date query string true "Start date (MM-YYYY)"
// @Param to_date query string true "End date (MM-YYYY)"
// @Param user_id query string false "Filter by user ID"
// @Param service_name query string false "Filter by service name"
// @Success 200 {object} models.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total-cost [get]
// func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
// 	var req models.TotalCostRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	total, err := h.service.GetTotalCost(c.Request.Context(), req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.TotalCostResponse{TotalCost: total})
// }
