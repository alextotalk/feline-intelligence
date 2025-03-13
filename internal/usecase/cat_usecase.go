package usecase

import (
	"context"
	"fmt"

	"github.com/alextotalk/feline-intelligence/internal/domain"
	"github.com/alextotalk/feline-intelligence/internal/domain/model"
	"github.com/alextotalk/feline-intelligence/internal/infrastructure/catapi"
)

type CatUsecase interface {
	CreateCat(ctx context.Context, cat *model.Cat) error
	GetCat(ctx context.Context, id int) (*model.Cat, error)
	ListCats(ctx context.Context) ([]model.Cat, error)
	UpdateCatSalary(ctx context.Context, catID int, newSalary float64) error
	DeleteCat(ctx context.Context, catID int) error
}

type catUsecase struct {
	catRepo domain.CatRepository
	catAPI  catapi.CatAPI // для валідації породи
}

func NewCatUsecase(cr domain.CatRepository, catAPI catapi.CatAPI) CatUsecase {
	return &catUsecase{
		catRepo: cr,
		catAPI:  catAPI,
	}
}

// CreateCat створює кота, перевіряючи, чи порода валідна (через TheCatAPI).
func (u *catUsecase) CreateCat(ctx context.Context, cat *model.Cat) error {
	valid, err := u.catAPI.IsBreedValid(ctx, cat.Breed)
	if err != nil {
		return fmt.Errorf("failed to validate cat breed: %w", err)
	}
	if !valid {
		return fmt.Errorf("breed '%s' is not a valid cat breed", cat.Breed)
	}
	return u.catRepo.Create(cat)
}

func (u *catUsecase) GetCat(ctx context.Context, id int) (*model.Cat, error) {
	return u.catRepo.GetByID(id)
}

func (u *catUsecase) ListCats(ctx context.Context) ([]model.Cat, error) {
	return u.catRepo.GetAll()
}

func (u *catUsecase) UpdateCatSalary(ctx context.Context, catID int, newSalary float64) error {
	cat, err := u.catRepo.GetByID(catID)
	if err != nil {
		return err
	}
	cat.Salary = newSalary
	return u.catRepo.Update(cat)
}

func (u *catUsecase) DeleteCat(ctx context.Context, catID int) error {
	return u.catRepo.Delete(catID)
}
