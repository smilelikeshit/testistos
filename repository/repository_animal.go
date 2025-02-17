package repository

import (
	"context"
	"database/sql"
	"fmt"
	"vault-app/domain"

	_ "github.com/lib/pq"
)

type repositoryAnimal struct {
	DB *sql.DB
}

func NewrepositoryAnimal(db *sql.DB) domain.AnimalRepository {
	return &repositoryAnimal{DB: db}
}

func (r *repositoryAnimal) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Animal, err error) {

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			return
		}
	}()

	result = make([]domain.Animal, 0)

	for rows.Next() {
		t := domain.Animal{}
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Age,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	fmt.Println(result)

	return result, nil
}

func (r *repositoryAnimal) Store(ctx context.Context, an *domain.Animal) (err error) {

	query := "INSERT INTO animals (name, age) VALUES ($1, $2) RETURNING id;"

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res := stmt.QueryRowContext(ctx, an.Name, an.Age).Scan(&an.Id)
	if err != nil {
		return
	}

	return res
}

func (r *repositoryAnimal) GetByID(ctx context.Context, id int) (res domain.Animal, err error) {

	query := "SELECT id, name, age FROM animals WHERE id = $1;"

	// err := r.DB.QueryRowContext(ctx, query, id).Scan(&animal.Id, &animal.Name, &animal.Age)
	// if err != nil {
	// 	return &animal, err
	// }

	list, err := r.fetch(ctx, query, id)

	if err != nil {
		return domain.Animal{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return res, err
}
