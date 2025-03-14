package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/alextotalk/feline-intelligence/internal/domain/model"
	"github.com/alextotalk/feline-intelligence/internal/usecase"
)

type CatHandler struct {
	catUC usecase.CatUsecase
}

// NewCatHandler реєструє хендлери котів у роутері Echo.
func NewCatHandler(e *echo.Echo, catUC usecase.CatUsecase) {
	handler := &CatHandler{catUC: catUC}

	// Маршрути
	e.POST("/cats", handler.CreateCat)
	e.GET("/cats", handler.ListCats)
	e.GET("/cats/:id", handler.GetCatByID)
	e.PUT("/cats/:id/salary", handler.UpdateSalary)
	e.DELETE("/cats/:id", handler.DeleteCat)
}

// CreateCat створює нового кота.
// @Summary Створити кота
// @Description Створює нового шпигунського кота з наданими даними
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body model.Cat true "Дані кота"
// @Success 201 {object} model.Cat
// @Failure 400 {object} map[string]string "Невірний запит"
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /cats [post]
func (h *CatHandler) CreateCat(c echo.Context) error {
	var cat model.Cat
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.catUC.CreateCat(context.Background(), &cat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, cat)
}

// ListCats повертає список усіх котів.
// @Summary Список котів
// @Description Отримує список усіх шпигунських котів
// @Tags cats
// @Accept json
// @Produce json
// @Success 200 {array} model.Cat
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /cats [get]
func (h *CatHandler) ListCats(c echo.Context) error {
	cats, err := h.catUC.ListCats(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

// GetCatByID повертає кота за його ID.
// @Summary Отримати кота за ID
// @Description Отримує деталі кота за його унікальним ідентифікатором
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID кота"
// @Success 200 {object} model.Cat
// @Failure 404 {object} map[string]string "Кіт не знайдений"
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /cats/{id} [get]
func (h *CatHandler) GetCatByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cat, err := h.catUC.GetCat(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if cat == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "cat not found"})
	}
	return c.JSON(http.StatusOK, cat)
}

// UpdateSalary оновлює зарплату кота.
// @Summary Оновити зарплату кота
// @Description Оновлює зарплату кота за його ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID кота"
// @Param salary body float64 true "Нова зарплата"
// @Success 200
// @Failure 400 {object} map[string]string "Невірний запит"
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /cats/{id}/salary [put]
func (h *CatHandler) UpdateSalary(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	type salaryReq struct {
		Salary float64 `json:"salary"`
	}
	var req salaryReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.catUC.UpdateCatSalary(context.Background(), id, req.Salary); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// DeleteCat видаляє кота за його ID.
// @Summary Видалити кота
// @Description Видаляє кота за його унікальним ідентифікатором
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID кота"
// @Success 200
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCat(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.catUC.DeleteCat(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
