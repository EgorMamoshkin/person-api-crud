package logic

import (
	"context"
	"fmt"
	"github.com/EgorMamoshkin/person-api-crud/internal/app"
	"time"
)

type PerLogic struct {
	perRepo    app.PersonRepository
	ctxTimeout time.Duration
}

func NewPersonLogic(perRep app.PersonRepository, timeout time.Duration) *PerLogic {
	return &PerLogic{perRepo: perRep, ctxTimeout: timeout}
}

func (p *PerLogic) StorePerson(ctx context.Context, per *app.Person) error {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	ok, err := p.isEmailExist(ctx, per.Email, 0)
	if err != nil {
		return fmt.Errorf("can't check is person already exist: %w", err)
	}

	if ok {
		return fmt.Errorf("another person with email address: %s already exist", per.Email)
	}

	return p.perRepo.Store(ctx, per)
}

func (p *PerLogic) DeletePerson(ctx context.Context, id int) error { //TODO
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	return p.perRepo.Delete(ctx, id)
}

func (p *PerLogic) GetPersonByID(ctx context.Context, id int) (*app.Person, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	person, err := p.perRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *PerLogic) UpdatePerson(ctx context.Context, per *app.Person) error {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	if err := p.isPersonExist(ctx, per.Id); err != nil {
		return fmt.Errorf("can't update person: %w", err)
	}

	ok, err := p.isEmailExist(ctx, per.Email, per.Id)
	if err != nil {
		return fmt.Errorf("can't check if the email is already using: %w", err)
	}

	if ok {
		return fmt.Errorf("another person already using this email address: %s", per.Email)
	}

	return p.perRepo.Update(ctx, per)
}

func (p *PerLogic) isEmailExist(ctx context.Context, email string, id int) (bool, error) { //TODO
	per, err := p.perRepo.GetByEmail(ctx, email, id)
	if err != nil {
		return false, err
	}

	if *per != (app.Person{}) {
		return true, nil
	}

	return false, nil
}

func (p *PerLogic) isPersonExist(ctx context.Context, id int) error {
	_, err := p.perRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
