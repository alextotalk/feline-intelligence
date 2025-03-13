package domain

import "github.com/alextotalk/feline-intelligence/internal/domain/model"

// CatRepository описує всі операції з котами (CRUD).
type CatRepository interface {
	Create(cat *model.Cat) error
	GetByID(id int) (*model.Cat, error)
	GetAll() ([]model.Cat, error)
	Update(cat *model.Cat) error
	Delete(id int) error
}

// MissionRepository описує CRUD-операції з місіями та додатковий метод призначення кота.
type MissionRepository interface {
	Create(mission *model.Mission) error
	GetByID(id int) (*model.Mission, error)
	GetAll() ([]model.Mission, error)
	Update(mission *model.Mission) error
	Delete(id int) error
	AssignCat(missionID, catID int) error
}

// TargetRepository описує операції з “цілями” (Target).
type TargetRepository interface {
	AddToMission(target *model.Target) error
	Update(target *model.Target) error
	Delete(id int) error
	GetByID(id int) (*model.Target, error) // За потреби
}
