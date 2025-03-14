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

func NewCatHandler(e *echo.Echo, catUC usecase.CatUsecase) {
	handler := &CatHandler{catUC: catUC}

	e.POST("/cats", handler.CreateCat)
	e.GET("/cats", handler.ListCats)
	e.GET("/cats/:id", handler.GetCatByID)
	e.PUT("/cats/:id/salary", handler.UpdateSalary)
	e.DELETE("/cats/:id", handler.DeleteCat)
}

// CreateCat Creates a new cat.
// @Summary Create a cat
// @Description Creates a new spy cat with data provided
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body model.Cat true "Cat data"
// @Success 201 {object} model.Cat
// @Failure 400 {object} map[string]string “Incorrect request”
// @Failure 500 {object} map[string]string "Internal server error"
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

// ListCats Returns all cats.
// @Summary List of cats
// @Description Gets a list of all spy cats
// @Tags cats
// @Accept json
// @Produce json
// @Success 200 {array} model.Cat
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cats [get]
func (h *CatHandler) ListCats(c echo.Context) error {
	cats, err := h.catUC.ListCats(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

// GetCatByID Returns the cat for his ID.
// @Summary Get a cat for ID
// @Description Receives cat details by its unique ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID cat"
// @Success 200 {object} model.Cat
// @Failure 404 {object} map[string]string "Cat is not found"
// @Failure 500 {object} map[string]string "Internal server error"
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

// UpdateSalary Updates a cat's salary.
// @Summary Update a cat's salary
// @Description Updates a cat's salary for his ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID кота"
// @Param salary body float64 true "New salary"
// @Success 200
// @Failure 400 {object} map[string]string “Incorrect request”
// @Failure 500 {object} map[string]string "Internal server error"
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

// DeleteCat Removes the cat for his ID.
// @Summary Remove the cat
// @Description Removes the cat by its unique ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "ID cat"
// @Success 200
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCat(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.catUC.DeleteCat(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
