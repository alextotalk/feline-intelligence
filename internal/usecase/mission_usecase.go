package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alextotalk/feline-intelligence/internal/domain"
	"github.com/alextotalk/feline-intelligence/internal/domain/model"
)

type MissionUsecase interface {
	CreateMission(ctx context.Context, mission *model.Mission) error
	DeleteMission(ctx context.Context, missionID int) error
	CompleteMission(ctx context.Context, missionID int) error

	GetMission(ctx context.Context, id int) (*model.Mission, error)
	ListMissions(ctx context.Context) ([]model.Mission, error)
	AssignCatToMission(ctx context.Context, missionID, catID int) error

	AddTarget(ctx context.Context, target *model.Target) error
	DeleteTarget(ctx context.Context, targetID int) error
	CompleteTarget(ctx context.Context, targetID int) error
	UpdateTargetNotes(ctx context.Context, targetID int, newNotes string) error
}

type missionUsecase struct {
	missionRepo domain.MissionRepository
	targetRepo  domain.TargetRepository
	catRepo     domain.CatRepository
}

func NewMissionUsecase(mr domain.MissionRepository, tr domain.TargetRepository, cr domain.CatRepository) MissionUsecase {
	return &missionUsecase{
		missionRepo: mr,
		targetRepo:  tr,
		catRepo:     cr,
	}
}

func (u *missionUsecase) CreateMission(ctx context.Context, mission *model.Mission) error {
	if err := u.missionRepo.Create(mission); err != nil {
		return err
	}
	for i := range mission.Targets {
		mission.Targets[i].MissionID = mission.ID
		if err := u.targetRepo.AddToMission(&mission.Targets[i]); err != nil {
			return err
		}
	}
	return nil
}

func (u *missionUsecase) DeleteMission(ctx context.Context, missionID int) error {
	mission, err := u.missionRepo.GetByID(missionID)
	if err != nil {
		return err
	}
	if mission.CatID != nil {
		return fmt.Errorf("cannot delete mission %d: it is assigned to cat", missionID)
	}
	return u.missionRepo.Delete(missionID)
}

func (u *missionUsecase) CompleteMission(ctx context.Context, missionID int) error {
	mission, err := u.missionRepo.GetByID(missionID)
	if err != nil {
		return err
	}
	// check that all goals are completed
	for _, t := range mission.Targets {
		if !t.Complete {
			return fmt.Errorf("target %d is not complete, cannot complete mission %d", t.ID, missionID)
		}
	}
	mission.Completed = true
	return u.missionRepo.Update(mission)
}

func (u *missionUsecase) GetMission(ctx context.Context, id int) (*model.Mission, error) {
	return u.missionRepo.GetByID(id)
}

func (u *missionUsecase) ListMissions(ctx context.Context) ([]model.Mission, error) {
	return u.missionRepo.GetAll()
}

func (u *missionUsecase) AssignCatToMission(ctx context.Context, missionID, catID int) error {
	// check if the cat exists
	cat, err := u.catRepo.GetByID(catID)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("cat does not exist")
	}
	// check if the mission is completed
	mission, err := u.missionRepo.GetByID(missionID)
	if err != nil {
		return err
	}
	if mission.Completed {
		return errors.New("cannot assign cat to a completed mission")
	}
	return u.missionRepo.AssignCat(missionID, catID)
}

func (u *missionUsecase) AddTarget(ctx context.Context, target *model.Target) error {
	mission, err := u.missionRepo.GetByID(target.MissionID)
	if err != nil {
		return err
	}
	if mission.Completed {
		return fmt.Errorf("cannot add target: mission %d is completed", mission.ID)
	}
	return u.targetRepo.AddToMission(target)
}

func (u *missionUsecase) DeleteTarget(ctx context.Context, targetID int) error {
	return u.targetRepo.Delete(targetID)
}

func (u *missionUsecase) CompleteTarget(ctx context.Context, targetID int) error {
	t, err := u.targetRepo.GetByID(targetID)
	if err != nil {
		return err
	}
	t.Complete = true
	return u.targetRepo.Update(t)
}

func (u *missionUsecase) UpdateTargetNotes(ctx context.Context, targetID int, newNotes string) error {
	t, err := u.targetRepo.GetByID(targetID)
	if err != nil {
		return err
	}
	// In the database, triggers check whether it is possible to update Notes.
	// We can additionally check at the business logic level:
	if t.Complete {
		return errors.New("cannot update notes of a completed target")
	}
	t.Notes = newNotes
	return u.targetRepo.Update(t)
}
