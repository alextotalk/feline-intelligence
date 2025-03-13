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

func (h *CatHandler) ListCats(c echo.Context) error {
	cats, err := h.catUC.ListCats(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

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

func (h *CatHandler) DeleteCat(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.catUC.DeleteCat(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
