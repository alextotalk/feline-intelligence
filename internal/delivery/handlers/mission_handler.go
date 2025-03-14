package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/alextotalk/feline-intelligence/internal/domain/model"
	"github.com/alextotalk/feline-intelligence/internal/usecase"
)

type MissionHandler struct {
	missionUC usecase.MissionUsecase
}

func NewMissionHandler(e *echo.Echo, mu usecase.MissionUsecase) {
	handler := &MissionHandler{missionUC: mu}

	// Місії
	e.POST("/missions", handler.CreateMission)
	e.GET("/missions", handler.ListMissions)
	e.GET("/missions/:id", handler.GetMission)
	e.PUT("/missions/:id/complete", handler.CompleteMission)
	e.DELETE("/missions/:id", handler.DeleteMission)

	// Призначити кота
	e.POST("/missions/:id/assign/:catID", handler.AssignCat)

	// Цілі
	e.POST("/missions/:id/targets", handler.AddTarget)
	e.DELETE("/targets/:targetID", handler.DeleteTarget)
	e.PUT("/targets/:targetID/complete", handler.CompleteTarget)
	e.PUT("/targets/:targetID/notes", handler.UpdateTargetNotes)
}

// CreateMission створює нову місію.
// @Summary Створити місію
// @Description Створює нову місію з наданими даними
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body model.Mission true "Дані місії"
// @Success 201 {object} model.Mission
// @Failure 400 {object} map[string]string "Невірний запит"
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /missions [post]
func (h *MissionHandler) CreateMission(c echo.Context) error {
	var mission model.Mission
	if err := c.Bind(&mission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.missionUC.CreateMission(context.Background(), &mission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, mission)
}

// ListMissions повертає список усіх місій.
// @Summary Список місій
// @Description Отримує список усіх місій
// @Tags missions
// @Accept json
// @Produce json
// @Success 200 {array} model.Mission
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /missions [get]
func (h *MissionHandler) ListMissions(c echo.Context) error {
	missions, err := h.missionUC.ListMissions(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, missions)
}

// GetMission повертає місію за її ID.
// @Summary Отримати місію за ID
// @Description Отримує деталі місії за її унікальним ідентифікатором
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Success 200 {object} model.Mission
// @Failure 404 {object} map[string]string "Місія не знайдена"
// @Failure 500 {object} map[string]string "Внутрішня помилка сервера"
// @Router /missions/{id} [get]
func (h *MissionHandler) GetMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	mission, err := h.missionUC.GetMission(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if mission == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "mission not found"})
	}
	return c.JSON(http.StatusOK, mission)
}

// CompleteMission завершує місію.
// @Summary Завершити місію
// @Description Позначає місію як завершену
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Success 200
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, місія вже завершена)"
// @Router /missions/{id}/complete [put]
func (h *MissionHandler) CompleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.CompleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// DeleteMission видаляє місію.
// @Summary Видалити місію
// @Description Видаляє місію за її ID
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Success 200
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, місія призначена коту)"
// @Router /missions/{id} [delete]
func (h *MissionHandler) DeleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.DeleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// AssignCat призначає кота до місії.
// @Summary Призначити кота до місії
// @Description Призначає кота до конкретної місії
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Param catID path int true "ID кота"
// @Success 200
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, кіт вже має активну місію)"
// @Router /missions/{id}/assign/{catID} [post]
func (h *MissionHandler) AssignCat(c echo.Context) error {
	missionID, _ := strconv.Atoi(c.Param("id"))
	catID, _ := strconv.Atoi(c.Param("catID"))
	if err := h.missionUC.AssignCatToMission(context.Background(), missionID, catID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// AddTarget додає нову ціль до місії.
// @Summary Додати ціль до місії
// @Description Додає нову ціль до конкретної місії
// @Tags targets
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Param target body model.Target true "Дані цілі"
// @Success 201 {object} model.Target
// @Failure 400 {object} map[string]string "Невірний запит"
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, місія завершена або має максимум цілей)"
// @Router /missions/{id}/targets [post]
func (h *MissionHandler) AddTarget(c echo.Context) error {
	missionID, _ := strconv.Atoi(c.Param("id"))
	var target model.Target
	if err := c.Bind(&target); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	target.MissionID = missionID
	if err := h.missionUC.AddTarget(context.Background(), &target); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, target)
}

// DeleteTarget видаляє ціль.
// @Summary Видалити ціль
// @Description Видаляє ціль за її ID
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID цілі"
// @Success 200
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, ціль завершена або місія має лише одну ціль)"
// @Router /targets/{targetID} [delete]
func (h *MissionHandler) DeleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.DeleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// CompleteTarget завершує ціль.
// @Summary Завершити ціль
// @Description Позначає ціль як завершену
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID цілі"
// @Success 200
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, ціль вже завершена)"
// @Router /targets/{targetID}/complete [put]
func (h *MissionHandler) CompleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.CompleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// UpdateTargetNotes оновлює нотатки цілі.
// @Summary Оновити нотатки цілі
// @Description Оновлює нотатки для конкретної цілі
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID цілі"
// @Param notes body string true "Нові нотатки"
// @Success 200
// @Failure 400 {object} map[string]string "Невірний запит"
// @Failure 409 {object} map[string]string "Конфлікт (наприклад, ціль або місія завершені)"
// @Router /targets/{targetID}/notes [put]
func (h *MissionHandler) UpdateTargetNotes(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	type notesReq struct {
		Notes string `json:"notes"`
	}
	var nr notesReq
	if err := c.Bind(&nr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.missionUC.UpdateTargetNotes(context.Background(), targetID, nr.Notes); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
