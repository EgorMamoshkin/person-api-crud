package app

import "context"

type Person struct {
	Id        int    `json:"id"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type PersonRepository interface {
	Store(ctx context.Context, person *Person) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Person, error)
	GetByEmail(ctx context.Context, email string, id int) (*Person, error)
	Update(ctx context.Context, person *Person) error
}

type PersonLogic interface {
	StorePerson(ctx context.Context, per *Person) error
	DeletePerson(ctx context.Context, id int) error
	GetPersonByID(ctx context.Context, id int) (*Person, error)
	UpdatePerson(ctx context.Context, per *Person) error
}
