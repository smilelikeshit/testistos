package usecase

import (
	"context"
	"vault-app/domain"
)

type UsecaseAnimal struct {
	IAnimalRepo domain.AnimalRepository
}

func NewusecaseAnimal(animalRepo *domain.AnimalRepository) domain.AnimalUseCase {
	return &UsecaseAnimal{
		IAnimalRepo: *animalRepo,
	}
}

func (u *UsecaseAnimal) Store(ctx context.Context, an *domain.Animal) (err error) {

	err = u.IAnimalRepo.Store(ctx, an)

	if err != nil {
		return err
	}
	return

}

func (u *UsecaseAnimal) GetByID(ctx context.Context, id int) (*domain.Animal, error) {

	animal, err := u.IAnimalRepo.GetByID(ctx, id)

	if err != nil {
		return &animal, err
	}

	return &animal, err

}
