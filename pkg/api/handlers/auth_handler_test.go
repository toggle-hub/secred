package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xsadia/secred/pkg/api/handlers"
	"github.com/xsadia/secred/pkg/database"
	usermodel "github.com/xsadia/secred/pkg/models/user_model"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	server *echo.Echo
	db     *sql.DB
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	godotenv.Load("../../../.env")
	storage, err := database.NewDB(
		"localhost",
		"root",
		"root",
		"secred_test",
		"disable",
		"file://../../../migrations",
	)
	if err != nil {
		log.Fatalf("failed to start test database %v", err.Error())
	}

	suite.db = storage.DB()
	suite.server = echo.New()

	suite.server.POST("/register", handlers.RegisterHandler)
	suite.server.POST("/login", handlers.LoginHandler)
}

func (suite *AuthHandlerTestSuite) AfterTest(_, _ string) {
	_, err := suite.db.Exec("TRUNCATE TABLE users, items, orders, order_items, schools, school_orders, school_items, school_order_items RESTART IDENTITY;")
	if err != nil {
		log.Fatalf("Failed to truncate %v", err.Error())
	}
}

func (suite *AuthHandlerTestSuite) TestRegisterSuccess() {
	t := suite.T()
	requestBody := []byte(`{
		"email": "fizi@gmail.com",
		"name": "fizi",
		"password": "123123123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	suite.server.ServeHTTP(recorder, request)
	var response handlers.AuthResponseBody
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	_, err := usermodel.New(suite.db).FindByEmail(response.User.Email)
	assert.NoError(t, err)
	assert.Equal(t, "fizi@gmail.com", response.User.Email)
}
func (suite *AuthHandlerTestSuite) TestRegisterConflict() {
	t := suite.T()
	userModel := usermodel.New(suite.db)
	userModel.Create("fizi@gmail.com", "fizi", "123123123")

	requestBody := []byte(`{
		"email": "fizi@gmail.com",
		"name": "fizi1",
		"password": "123123123"
	}`)

	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	suite.server.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusConflict, recorder.Code)
}

func (suite *AuthHandlerTestSuite) TestLoginSuccess() {
	t := suite.T()
	requestBody := []byte(`{
		"email": "fizi@gmail.com",
		"password": "123123123"
	}`)

	userModel := usermodel.New(suite.db)
	userModel.Create("fizi@gmail.com", "fizi", "123123123")

	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	suite.server.ServeHTTP(recorder, request)
	var response handlers.AuthResponseBody
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	_, err := userModel.FindByEmail(response.User.Email)
	assert.NoError(t, err)
	assert.Equal(t, "fizi@gmail.com", response.User.Email)
}

func (suite *AuthHandlerTestSuite) TestLoginWrongPassword() {
	t := suite.T()
	requestBody := []byte(`{
		"email": "fizi@gmail.com",
		"password": "123123122"
	}`)

	userModel := usermodel.New(suite.db)
	userModel.Create("fizi@gmail.com", "fizi", "123123123")

	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	suite.server.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
