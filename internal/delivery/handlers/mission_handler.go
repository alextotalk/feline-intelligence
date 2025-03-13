package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alextotalk/feline-intelligence/internal/domain/model"
	"github.com/alextotalk/feline-intelligence/internal/usecase"
	"github.com/labstack/echo/v4"
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

// Створити місію
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

func (h *MissionHandler) ListMissions(c echo.Context) error {
	missions, err := h.missionUC.ListMissions(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, missions)
}

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

func (h *MissionHandler) CompleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.CompleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *MissionHandler) DeleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.DeleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// Призначити кота до місії
func (h *MissionHandler) AssignCat(c echo.Context) error {
	missionID, _ := strconv.Atoi(c.Param("id"))
	catID, _ := strconv.Atoi(c.Param("catID"))
	if err := h.missionUC.AssignCatToMission(context.Background(), missionID, catID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// Цілі
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

func (h *MissionHandler) DeleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.DeleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *MissionHandler) CompleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.CompleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

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
