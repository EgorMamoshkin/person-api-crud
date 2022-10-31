package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/EgorMamoshkin/person-api-crud/internal/app"
	"github.com/EgorMamoshkin/person-api-crud/internal/app/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPersonHandler_StorePerson(t *testing.T) {
	type mockBehavior func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person)

	testTable := []struct {
		name      string
		inputBody string
		inputUser *app.Person
		mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@gmail.com", "phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			inputUser: &app.Person{
				Id:        0,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {
				s.EXPECT().StorePerson(ctx, pers).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":0,"email":"test@gmail.com","phone":"+1111111111","firstName":"TestName","lastName":"Test"}`,
		}, {
			name:                "Empty Fields",
			inputBody:           `{"phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {},
			expectedStatusCode:  400,
			expectedRequestBody: `"invalid request data: Key: 'Person.Email' Error:Field validation for 'Email' failed on the 'required' tag"`,
		}, {
			name:      "Service Failure",
			inputBody: `{"id":1, "email":"test@gmail.com", "phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			inputUser: &app.Person{
				Id:        1,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {
				s.EXPECT().StorePerson(context.Background(), pers).Return(errors.New("service failure"))
			},
			expectedStatusCode:  501,
			expectedRequestBody: `"service failure"`,
		}, {
			name:      "Unprocessable Entity",
			inputBody: `111111111`,
			inputUser: &app.Person{
				Id:        1,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {},
			expectedStatusCode:  422,
			expectedRequestBody: `"code=400, message=Unmarshal type error: expected=app.Person, got=number, field=, offset=9, internal=json: cannot unmarshal number into Go value of type app.Person"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			perLog := mock_app.NewMockPersonLogic(ctrl)
			testCase.mockBehavior(perLog, context.Background(), testCase.inputUser)

			hand := PersonHandler{perLog}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/person", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			r := echo.New()
			r.POST("/person", hand.StorePerson)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, testCase.expectedStatusCode)
			require.Equal(t, testCase.expectedRequestBody, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}
}

func TestPersonHandler_GetPerson(t *testing.T) {
	type mockBehavior func(s *mock_app.MockPersonLogic, ctx context.Context, id any)

	testTable := []struct {
		name    string
		inputID any
		mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {
				s.EXPECT().GetPersonByID(ctx, id).Return(&app.Person{Id: 1, Email: "test@gmail.com", Phone: "+1111111", FirstName: "Test", LastName: "Test"}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"email":"test@gmail.com","phone":"+1111111","firstName":"Test","lastName":"Test"}`,
		}, {
			name:                "Wrong ID",
			inputID:             "a",
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {},
			expectedStatusCode:  400,
			expectedRequestBody: `"strconv.Atoi: parsing \"a\": invalid syntax"`,
		}, {
			name:    "Service Failure",
			inputID: 0,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {
				s.EXPECT().GetPersonByID(ctx, id).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:  501,
			expectedRequestBody: `"service failure"`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			perLog := mock_app.NewMockPersonLogic(ctrl)
			testCase.mockBehavior(perLog, context.Background(), testCase.inputID)

			hand := PersonHandler{perLog}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/person/%v", testCase.inputID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			r := echo.New()
			r.GET("/person/:id", hand.GetPerson)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, testCase.expectedStatusCode)
			require.Equal(t, testCase.expectedRequestBody, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}

}

func TestPersonHandler_DeletePerson(t *testing.T) {
	type mockBehavior func(s *mock_app.MockPersonLogic, ctx context.Context, id any)

	testTable := []struct {
		name    string
		inputID any
		mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "OK",
			inputID: 0,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {
				s.EXPECT().DeletePerson(ctx, id).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `"The person's data has been deleted"`,
		}, {
			name:                "Wrong ID",
			inputID:             "a",
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {},
			expectedStatusCode:  400,
			expectedRequestBody: `"strconv.Atoi: parsing \"a\": invalid syntax"`,
		}, {
			name:    "Service Failure",
			inputID: 0,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, id any) {
				s.EXPECT().DeletePerson(ctx, id).Return(errors.New("service failure"))
			},
			expectedStatusCode:  501,
			expectedRequestBody: `"service failure"`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			perLog := mock_app.NewMockPersonLogic(ctrl)
			testCase.mockBehavior(perLog, context.Background(), testCase.inputID)

			hand := PersonHandler{perLog}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/person/%v", testCase.inputID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			r := echo.New()
			r.DELETE("/person/:id", hand.DeletePerson)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, testCase.expectedStatusCode)
			require.Equal(t, testCase.expectedRequestBody, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}

}

func TestPersonHandler_UpdatePerson(t *testing.T) {
	type mockBehavior func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person)

	testTable := []struct {
		name      string
		inputBody string
		inputUser *app.Person
		mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1, "email":"test@gmail.com", "phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			inputUser: &app.Person{
				Id:        1,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {
				s.EXPECT().UpdatePerson(ctx, pers).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"email":"test@gmail.com","phone":"+1111111111","firstName":"TestName","lastName":"Test"}`,
		}, {
			name:                "Empty Fields",
			inputBody:           `{"phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {},
			expectedStatusCode:  400,
			expectedRequestBody: `"invalid request data: Key: 'Person.Email' Error:Field validation for 'Email' failed on the 'required' tag"`,
		}, {
			name:      "Service Failure",
			inputBody: `{"id":1, "email":"test@gmail.com", "phone":"+1111111111", "firstName":"TestName", "lastname":"Test"}`,
			inputUser: &app.Person{
				Id:        1,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {
				s.EXPECT().UpdatePerson(context.Background(), pers).Return(errors.New("service failure"))
			},
			expectedStatusCode:  501,
			expectedRequestBody: `"service failure"`,
		}, {
			name:      "Unprocessable Entity",
			inputBody: `111111111`,
			inputUser: &app.Person{
				Id:        1,
				Email:     "test@gmail.com",
				Phone:     "+1111111111",
				FirstName: "TestName",
				LastName:  "Test",
			},
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, pers *app.Person) {},
			expectedStatusCode:  422,
			expectedRequestBody: `"code=400, message=Unmarshal type error: expected=app.Person, got=number, field=, offset=9, internal=json: cannot unmarshal number into Go value of type app.Person"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			perLog := mock_app.NewMockPersonLogic(ctrl)
			testCase.mockBehavior(perLog, context.Background(), testCase.inputUser)

			hand := PersonHandler{perLog}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/person", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			r := echo.New()
			r.PUT("/person", hand.UpdatePerson)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, testCase.expectedStatusCode)
			require.Equal(t, testCase.expectedRequestBody, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}
}

func TestPersonHandler_GetPersonList(t *testing.T) {
	type mockBehavior func(s *mock_app.MockPersonLogic, ctx context.Context, offsetId int, batchSize any)

	testTable := []struct {
		name       string
		inputID    int
		inputBSize any
		mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			inputID:    0,
			inputBSize: 2,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, offsetId int, batchSize any) {
				s.EXPECT().GetPersonList(ctx, offsetId, batchSize).Return([]app.Person{{1, "test@gmail.com", "+111111", "Test", "Test"}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":1,"email":"test@gmail.com","phone":"+111111","firstName":"Test","lastName":"Test"}]`,
		}, {
			name:                "Wrong ID",
			inputID:             0,
			inputBSize:          "a",
			mockBehavior:        func(s *mock_app.MockPersonLogic, ctx context.Context, offsetId int, batchSize any) {},
			expectedStatusCode:  400,
			expectedRequestBody: `"strconv.Atoi: parsing \"a\": invalid syntax"`,
		}, {
			name:       "Service Failure",
			inputID:    0,
			inputBSize: 2,
			mockBehavior: func(s *mock_app.MockPersonLogic, ctx context.Context, offsetId int, batchSize any) {
				s.EXPECT().GetPersonList(ctx, offsetId, batchSize).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:  501,
			expectedRequestBody: `"service failure"`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			perLog := mock_app.NewMockPersonLogic(ctrl)
			testCase.mockBehavior(perLog, context.Background(), testCase.inputID, testCase.inputBSize)

			hand := PersonHandler{perLog}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/person/%v/%v", testCase.inputID, testCase.inputBSize), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			r := echo.New()
			r.GET("/person/:offsetId/:batchSize", hand.GetPersonList)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, testCase.expectedStatusCode)
			require.Equal(t, testCase.expectedRequestBody, strings.TrimRight(rec.Body.String(), "\n"))
		})
	}

}
