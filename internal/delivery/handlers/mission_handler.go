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

	e.POST("/missions", handler.CreateMission)
	e.GET("/missions", handler.ListMissions)
	e.GET("/missions/:id", handler.GetMission)
	e.PUT("/missions/:id/complete", handler.CompleteMission)
	e.DELETE("/missions/:id", handler.DeleteMission)

	e.POST("/missions/:id/assign/:catID", handler.AssignCat)

	e.POST("/missions/:id/targets", handler.AddTarget)
	e.DELETE("/targets/:targetID", handler.DeleteTarget)
	e.PUT("/targets/:targetID/complete", handler.CompleteTarget)
	e.PUT("/targets/:targetID/notes", handler.UpdateTargetNotes)
}

// CreateMission Creates a new mission.
// @Summary Create a mission
// @Description Creates a new mission with the data provided
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body model.Mission true "Mission data"
// @Success 201 {object} model.Mission
// @Failure 400 {object} map[string]string “Incorrect request”
// @Failure 500 {object} map[string]string "Internal server error"
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

// ListMissions returns a list of all missions.
// @Summary List of missions
// @Description Gets a list of all missions
// @Tags missions
// @Accept json
// @Produce json
// @Success 200 {array} model.Mission
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /missions [get]
func (h *MissionHandler) ListMissions(c echo.Context) error {
	missions, err := h.missionUC.ListMissions(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, missions)
}

// GetMission Returns the mission for her ID.
// @Summary Get a mission for ID
// @Description Receives mission details by its unique identifier
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID mission"
// @Success 200 {object} model.Mission
// @Failure 404 {object} map[string]string "Mission not found"
// @Failure 500 {object} map[string]string "Internal server error"
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

// CompleteMission completes the mission.
// @Summary Complete the mission
// @Description Denotes the mission as completed
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID місії"
// @Success 200
// @Failure 409 {object} map[string]string “Conflict (mission is already completed)”
// @Router /missions/{id}/complete [put]
func (h *MissionHandler) CompleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.CompleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// DeleteMission Removes the mission.
// @Summary Remove the mission
// @Description Removes the mission for her ID
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID mission"
// @Success 200
// @Failure 409 {object} map[string]string "Conflict ( mission assigned a cat)"
// @Router /missions/{id} [delete]
func (h *MissionHandler) DeleteMission(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.missionUC.DeleteMission(context.Background(), id); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// AssignCat appoints a cat to a mission.
// @Summary To assign a cat to a mission
// @Description Appoints a cat to a particular mission
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "ID mission"
// @Param catID path int true "ID cat"
// @Success 200
// @Failure 409 {object} map[string]string "Conflict ( cat already has an active mission)"
// @Router /missions/{id}/assign/{catID} [post]
func (h *MissionHandler) AssignCat(c echo.Context) error {
	missionID, _ := strconv.Atoi(c.Param("id"))
	catID, _ := strconv.Atoi(c.Param("catID"))
	if err := h.missionUC.AssignCatToMission(context.Background(), missionID, catID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// AddTarget adds a new target to the mission.
// @Summary Add the target to the mission
// @Description Adds a new target to a particular mission
// @Tags targets
// @Accept json
// @Produce json
// @Param id path int true "ID mission"
// @Param target body model.Target true "Дані цілі"
// @Success 201 {object} model.Target
// @Failure 400 {object} map[string]string “Incorrect request”
// @Failure 409 {object} map[string]string "Conflict (mission is completed or has a maximum of purposes)"
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

// DeleteTarget Removes the target.
// @Summary Remove the target
// @Description Removes the target for her ID
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID targets"
// @Success 200
// @Failure 409 {object} map[string]string "Conflict (eg target completed or mission has only one target)"
// @Router /targets/{targetID} [delete]
func (h *MissionHandler) DeleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.DeleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// CompleteTarget completes the target.
// @Summary Complete the target
// @Description Denotes the target as completed
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID targets"
// @Success 200
// @Failure 409 {object} map[string]string “Conflict (eg target already completed)”
// @Router /targets/{targetID}/complete [put]
func (h *MissionHandler) CompleteTarget(c echo.Context) error {
	targetID, _ := strconv.Atoi(c.Param("targetID"))
	if err := h.missionUC.CompleteTarget(context.Background(), targetID); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

// UpdateTargetNotes Updates target notes.
// @Summary Update goals notes
// @Description Updates notes for a particular purpose
// @Tags targets
// @Accept json
// @Produce json
// @Param targetID path int true "ID targets"
// @Param notes body string true "New notes"
// @Success 200
// @Failure 400 {object} map[string]string “Incorrect request”
// @Failure 409 {object} map[string]string “Conflict (target or mission completed)”
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
