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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xsadia/secred/pkg/api/handlers"
	"github.com/xsadia/secred/pkg/database"
	"github.com/xsadia/secred/pkg/middlewares"
	schoolmodel "github.com/xsadia/secred/pkg/models/schools"
	usermodel "github.com/xsadia/secred/pkg/models/user_model"
	"github.com/xsadia/secred/pkg/utils"
)

type SchoolHandlerTestSuite struct {
	suite.Suite
	server *echo.Echo
	db     *sql.DB
}

func (suite *SchoolHandlerTestSuite) SetupTest() {
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
	suite.server.Use(middlewares.AuthMiddleware)
	suite.server.POST("/school", handlers.CreateSchoolHandler)
}

func (suite *SchoolHandlerTestSuite) AfterTest(_, _ string) {
	_, err := suite.db.Exec("TRUNCATE TABLE users, items, orders, order_items, schools, school_orders, school_items, school_order_items RESTART IDENTITY;")
	if err != nil {
		log.Fatalf("Failed to truncate %v", err.Error())
	}
}

func (suite *SchoolHandlerTestSuite) TestCreateSchoolSuccess() {
	t := suite.T()
	user, err := usermodel.New(suite.db).Create("fizi@gmail", "fizi", "123123123")
	assert.NoError(t, err)

	token, err := utils.CreateJWT(user.ID, utils.JWTExpireTime)
	assert.NoError(t, err)

	requestBody := []byte(`{
    "name": "CSC",
    "address": "rua vitor konder"
}`)

	request := httptest.NewRequest(http.MethodPost, "/school", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	recorder := httptest.NewRecorder()

	suite.server.ServeHTTP(recorder, request)
	var response schoolmodel.School
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))

	school, err := schoolmodel.New(suite.db).FindById(response.ID)
	assert.NoError(t, err)
	assert.Equal(t, school.Name, response.Name)
	assert.Equal(t, school.ID, response.ID)
}

func TestSchoolHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SchoolHandlerTestSuite))
}
