package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/EgorMamoshkin/person-api-crud/internal/app"
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type PSQLRepo struct {
	session *dbr.Session
}

func NewPostgresRepo(dsn string) *PSQLRepo {
	conn, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		logrus.Fatalf("failed to open a database : %s", err)
	}

	conn.SetMaxOpenConns(10)

	sess := conn.NewSession(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sess.PingContext(ctx); err != nil {
		logrus.Fatalf("can't connect to database: %s", err)
	}

	return &PSQLRepo{session: sess}
}

func (r *PSQLRepo) Store(ctx context.Context, person *app.Person) error {

	_, err := r.session.InsertInto("person").
		Columns("email", "phone", "first_name", "last_name").
		Record(person).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("can't save person: %w", err)
	}

	id, err := r.getID(ctx, person.Email)
	if err != nil {
		return err
	} else {
		person.Id = id
	}

	return nil
}

func (r *PSQLRepo) Delete(ctx context.Context, id int) error {
	res, err := r.session.DeleteFrom("person").Where("id = ?", id).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("can't delete person: %w", err)
	}
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't confirm is person deleted: %w", err)
	}

	if rowsDeleted == 0 {
		return errors.New("can't delete. person doesn't exist")
	}

	return nil
}

func (r *PSQLRepo) GetByID(ctx context.Context, id int) (*app.Person, error) {
	var person app.Person

	res, err := r.session.Select("*").From("person").
		Where("id = ?", id).LoadContext(ctx, &person)
	if err != nil {
		return nil, fmt.Errorf("can't get person: %w", err)
	}

	if res == 0 {
		return nil, fmt.Errorf("person with ID %d doesn't exist", id)
	}

	return &person, nil
}

func (r *PSQLRepo) GetByEmail(ctx context.Context, email string, id int) (*app.Person, error) {
	var person app.Person

	_, err := r.session.Select("*").From("person").
		Where("email = ? AND id <> ?", email, id).LoadContext(ctx, &person)

	if err != nil {
		return nil, fmt.Errorf("can't get person: %w", err)
	}

	return &person, nil
}

func (r *PSQLRepo) Update(ctx context.Context, per *app.Person) error {
	res, err := r.session.Update("person").
		Set("email", per.Email).
		Set("phone", per.Phone).
		Set("first_name", per.FirstName).
		Set("last_name", per.LastName).
		Where("id = ?", per.Id).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("can't update person: %w", err)
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		return fmt.Errorf("updating error. Affected rows: %d", numRows)
	}

	return nil
}

func (r *PSQLRepo) getID(ctx context.Context, email string) (int, error) {
	var id int

	res, err := r.session.Select("id").From("person").
		Where("email = ?", email).LoadContext(ctx, &id)

	if err != nil || res == 0 {
		return 0, fmt.Errorf("can't get ID: %w", err)
	}

	return id, nil
}
