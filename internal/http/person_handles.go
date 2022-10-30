package http

import (
	"fmt"
	"github.com/EgorMamoshkin/person-api-crud/internal/app"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	personLogic app.PersonLogic
}

func NewPersonHandler(e *echo.Echo, pu app.PersonLogic) {
	handler := &PersonHandler{personLogic: pu}

	e.GET("/person/:id", handler.GetPerson)
	e.POST("/person", handler.StorePerson)
	e.PUT("/person", handler.UpdatePerson)
	e.DELETE("/person/:id", handler.DeletePerson)

	log := logrus.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    v.URI,
				"status": v.Status,
				"ERROR":  v.Error,
				"Method": v.Method,
			}).Info("request")

			return nil
		},
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogMethod: true,
	}))
}

func (ph *PersonHandler) StorePerson(c echo.Context) error {
	var person app.Person

	err := c.Bind(&person)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&person); !ok {
		logrus.Error(err)

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	err = ph.personLogic.StorePerson(ctx, &person)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusNotImplemented, err.Error())
	}

	return c.JSON(http.StatusCreated, person)
}

func (ph *PersonHandler) GetPerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	person, err := ph.personLogic.GetPersonByID(ctx, id)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusNotImplemented, err.Error())
	}

	return c.JSON(http.StatusOK, *person)
}

func (ph *PersonHandler) DeletePerson(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	err = ph.personLogic.DeletePerson(ctx, id)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusNotImplemented, err.Error())
	}

	return c.JSON(http.StatusOK, "The person's data has been deleted")
}

func (ph *PersonHandler) UpdatePerson(c echo.Context) error {
	var person app.Person

	err := c.Bind(&person)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&person); !ok {
		logrus.Error(err)

		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	err = ph.personLogic.UpdatePerson(ctx, &person)
	if err != nil {
		logrus.Error(err)

		return c.JSON(http.StatusNotImplemented, err.Error())
	}

	return c.JSON(http.StatusOK, person)
}

func isRequestValid(p *app.Person) (bool, error) {
	val := validator.New()

	err := val.Struct(p)
	if err != nil {
		return false, fmt.Errorf("invalid request data: %w", err)
	}

	return true, nil
}
