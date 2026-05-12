package handlers

import (
	"errors"
	"net/http"

	subscrerrors "github.com/IvanKuchsh-600/subtracker/internal/domain/errors"
	subscrusecase "github.com/IvanKuchsh-600/subtracker/internal/usecase/subscription"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	usecase subscrusecase.Usecase
}

func NewSubscriptionHandler(usecase subscrusecase.Usecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: usecase}
}

// Create godoc
// @Summary Создание новой подписки
// @Description Добавляет новую подписку в систему
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} SubscriptionResponse "Подписка создана"
// @Failure 400 {object} ValidationErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	created, err := h.usecase.Create(c, toCreateInput(req))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toSubscriptionResponse(created))
}

// GetByID godoc
// @Summary Получение подписки по ID
// @Description Возвращает подписку по её уникальному идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} SubscriptionResponse "Подписка найдена"
// @Failure 400 {object} ValidationErrorResponse "Неверный ID"
// @Failure 404 {object} NotFoundErrorResponse "Подписка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	sub, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSubscriptionResponse(sub))
}

// Update godoc
// @Summary Обновление подписки
// @Description Обновляет поля существующей подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки" example("550e8400-e29b-41d4-a716-446655440000")
// @Param request body UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} SubscriptionResponse "Подписка обновлена"
// @Failure 400 {object} ValidationErrorResponse "Неверный запрос"
// @Failure 404 {object} NotFoundErrorResponse "Подписка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updated, err := h.usecase.Update(c, id, toUpdateInput(req))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSubscriptionResponse(updated))
}

// Delete godoc
// @Summary Удаление подписки
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 204 "Подписка удалена"
// @Failure 400 {object} ValidationErrorResponse "Неверный ID"
// @Failure 404 {object} NotFoundErrorResponse "Подписка не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	err := h.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary Список всех подписок
// @Description Возвращает список всех подписок в системе
// @Tags subscriptions
// @Produce json
// @Success 200 {array} SubscriptionResponse "Список подписок"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	subscriptions, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, toSubscriptionResponseList(subscriptions))
}

// GetTotalCost godoc
// @Summary Подсчёт стоимости подписок за период
// @Description Рассчитывает суммарную стоимость подписок, активных в указанном периоде
// @Tags subscriptions
// @Produce json
// @Param from_date query string true "Дата начала (MM-YYYY)" example("01-2025")
// @Param to_date query string true "Дата окончания (MM-YYYY)" example("12-2025")
// @Param user_id query string false "ID пользователя (UUID)" example("60601fee-2bf1-4721-ae6f-7636e79a0cba")
// @Param service_name query string false "Название сервиса" example("Yandex Plus")
// @Success 200 {object} TotalCostResponse "Общая стоимость"
// @Failure 400 {object} ValidationErrorResponse "Неверные параметры"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	var req TotalCostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	total, err := h.usecase.GetTotalCost(c.Request.Context(), toTotalCostRequest(req))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TotalCostResponse{TotalCost: total})
}

func (h *SubscriptionHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, subscrerrors.ErrNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "subscription not found"})
	case errors.Is(err, subscrerrors.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
	}
}
