package domain

import "github.com/alextotalk/feline-intelligence/internal/domain/model"

// CatRepository
type CatRepository interface {
	Create(cat *model.Cat) error
	GetByID(id int) (*model.Cat, error)
	GetAll() ([]model.Cat, error)
	Update(cat *model.Cat) error
	Delete(id int) error
}

// MissionRepository
type MissionRepository interface {
	Create(mission *model.Mission) error
	GetByID(id int) (*model.Mission, error)
	GetAll() ([]model.Mission, error)
	Update(mission *model.Mission) error
	Delete(id int) error
	AssignCat(missionID, catID int) error
}

// TargetRepository
type TargetRepository interface {
	AddToMission(target *model.Target) error
	Update(target *model.Target) error
	Delete(id int) error
	GetByID(id int) (*model.Target, error) // За потреби
}
