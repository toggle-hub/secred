package server_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xsadia/secred/graph"
	"github.com/xsadia/secred/graph/model"
	"github.com/xsadia/secred/pkg/database"
	itemmodel "github.com/xsadia/secred/pkg/models/item_model"
	schoolmodel "github.com/xsadia/secred/pkg/models/school_model"
	usermodel "github.com/xsadia/secred/pkg/models/user_model"
)

type Me struct {
	ID       string
	Name     string
	Email    string
	Typename string `json:"__typename"`
}
type CreateUser struct {
	Me Me
}

type CreateUserResponse struct {
	CreateUser CreateUser
}

type CreateSchool struct {
	ID          string
	Name        string
	Typename    string `json:"__typename"`
	Address     *string
	PhoneNumber *string
	Orders      []*model.Order
}

type CreateSchoolResponse struct {
	CreateSchool
}

type CreateItem struct {
	ID       string
	Name     string
	RawName  string
	Typename string `json:"__typename"`
	Quantity int
}

type CreateItemResponse struct {
	CreateItem
}

type MeResponse struct {
	Me
}

type GqlError struct {
	Message string
	Path    []string
}

type SecredTestSuite struct {
	suite.Suite
	client *client.Client
	db     *sql.DB
}

func addUserToContext(ctx context.Context, id string) client.Option {
	return func(r *client.Request) {
		ctx = context.WithValue(ctx, "user", id)
		r.HTTP = r.HTTP.WithContext(ctx)
	}
}

func (suite *SecredTestSuite) SetupTest() {
	db, err := database.New("localhost", "root", "root", "secred_test", "disable")
	if err != nil {
		log.Fatalf("failed to start test database %v", err.Error())
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate test database %v", err.Error())
	}

	suite.db = db
	suite.client = client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}})))
}

func (suite *SecredTestSuite) AfterTest(_, _ string) {
	_, err := suite.db.Exec("TRUNCATE TABLE users, items, orders, order_items, schools, school_orders, school_items, school_order_items RESTART IDENTITY;")
	if err != nil {
		log.Fatalf("Failed to truncate %v", err.Error())
	}
}

func (suite *SecredTestSuite) TestCreateSchoolSuccess() {
	t := suite.T()

	userModel := usermodel.New(suite.db)
	user, err := userModel.Create(model.CreateUserInput{Email: "fizi@gmail.com", Name: "fizi", Password: "123123"})
	require.NoError(t, err)

	var actual CreateSchoolResponse
	suite.client.MustPost(`mutation {
		createSchool(
			input: {name: "CSC", address: "R. Frei Evaristo, n 91"}
		) {
			id
			name
			address
			phoneNumber
			orders {
				id
			}
			__typename
		}
	}`, &actual, addUserToContext(context.Background(), user.ID))

	expectedAddr := "R. Frei Evaristo, n 91"
	expected := CreateSchoolResponse{
		CreateSchool: CreateSchool{
			ID:          actual.ID,
			Name:        "CSC",
			Address:     &expectedAddr,
			PhoneNumber: nil,
			Orders:      []*model.Order{},
			Typename:    "School",
		},
	}

	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateSchoolUnauthorized() {
	t := suite.T()

	resp, err := suite.client.RawPost(`mutation {
		createSchool(
			input: {name: "CSC", address: "R. Frei Evaristo, n 91"}
		) {
			id
			name
			address
			phoneNumber
			orders {
				id
			}
			__typename
		}
	}`)

	require.NoError(t, err)
	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "authorization required", Path: []string{"createSchool"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateSchoolDuplicateName() {
	t := suite.T()
	userModel := usermodel.New(suite.db)
	user, err := userModel.Create(model.CreateUserInput{
		Name:     "fizi",
		Email:    "fizi@gmail.com",
		Password: "123123",
	})
	require.NoError(t, err)

	schoolModel := schoolmodel.New(suite.db)
	_, err = schoolModel.Create(model.CreateSchoolInput{
		Name: "CSC",
	})
	require.NoError(t, err)

	resp, err := suite.client.RawPost(`mutation {
		createSchool(
			input: {name: "CSC", address: "R. Frei Evaristo, n 91"}
		) {
			id
			name
			address
			phoneNumber
			orders {
				id
			}
			__typename
		}
	}`, addUserToContext(context.Background(), user.ID))

	require.NoError(t, err)
	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "school name already in use", Path: []string{"createSchool"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateItemSuccess() {
	t := suite.T()

	userModel := usermodel.New(suite.db)
	user, err := userModel.Create(model.CreateUserInput{Email: "fizi@gmail.com", Name: "fizi", Password: "123123"})
	require.NoError(t, err)

	var actual CreateItemResponse
	suite.client.MustPost(`mutation {
		createItem(input: {name: "Feijão", quantity: 10}) {
			id
			name
			rawName
			quantity
			__typename
		}
	}`, &actual, addUserToContext(context.Background(), user.ID))

	expected := CreateItemResponse{
		CreateItem: CreateItem{
			ID:       actual.ID,
			Name:     "feijao",
			RawName:  "feijão",
			Quantity: 10,
			Typename: "Item",
		},
	}

	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateItemUnauthorized() {
	t := suite.T()

	resp, err := suite.client.RawPost(`mutation {
		createItem(input: {name: "Feijão", quantity: 10}) {
			id
			name
			rawName
			quantity
			__typename
		}
	}`)

	require.NoError(t, err)
	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "authorization required", Path: []string{"createItem"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateItemDuplicateName() {
	t := suite.T()
	userModel := usermodel.New(suite.db)
	user, err := userModel.Create(model.CreateUserInput{
		Name:     "fizi",
		Email:    "fizi@gmail.com",
		Password: "123123",
	})
	require.NoError(t, err)

	itemModel := itemmodel.New(suite.db)
	_, err = itemModel.Create(model.CreateItemInput{
		Name:     "féijão",
		Quantity: 10,
	})
	require.NoError(t, err)

	resp, err := suite.client.RawPost(`mutation {
		createItem(input: {name: "Feijão", quantity: 10}) {
			id
			name
			rawName
			quantity
			__typename
		}
	}`, addUserToContext(context.Background(), user.ID))

	require.NoError(t, err)
	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "item with given name already registered", Path: []string{"createItem"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateUserSuccess() {
	t := suite.T()

	var actual CreateUserResponse
	suite.client.MustPost(`mutation {
			createUser(input:{email:"fizi@gmail.com", name:"fizi", password:"123123"}) {
				me {
					name
					email
					__typename
				}
			}
		}`, &actual)

	expected := CreateUserResponse{
		CreateUser: CreateUser{
			Me: Me{
				Name:     "fizi",
				Email:    "fizi@gmail.com",
				Typename: "User",
			},
		},
	}

	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestCreateUserDuplicateEmail() {
	userModel := usermodel.New(suite.db)
	userModel.Create(model.CreateUserInput{
		Email:    "fizi@gmail.com",
		Name:     "fizi",
		Password: "123123",
	})

	t := suite.T()

	resp, err := suite.client.RawPost(`mutation {
			createUser(input:{email:"fizi@gmail.com", name:"fizi", password:"123123"}) {
				me {
					name
					email
					__typename
				}
			}
		}`)
	require.NoError(t, err)

	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "email already in use", Path: []string{"createUser"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestMeSuccess() {
	t := suite.T()
	userModel := usermodel.New(suite.db)
	createdUser, err := userModel.Create(model.CreateUserInput{
		Email:    "fizi@gmail.com",
		Name:     "fizi",
		Password: "123123",
	})

	require.NoError(t, err)

	var actual MeResponse
	suite.client.MustPost(`query {
				me {
					id
					name
					email
					__typename
				}
			}`, &actual, addUserToContext(context.Background(), createdUser.ID))

	expected := MeResponse{Me: Me{
		ID:       createdUser.ID,
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		Typename: "User",
	}}

	require.EqualValues(t, expected, actual)
}

func (suite *SecredTestSuite) TestMeUnauthorized() {
	t := suite.T()

	resp, err := suite.client.RawPost(`query {
				me {
					id
					name
					email
					__typename
				}
			}`)

	require.NoError(t, err)
	var actual []GqlError
	err = json.Unmarshal(resp.Errors, &actual)
	require.NoError(t, err)

	expected := []GqlError{
		{Message: "authorization required", Path: []string{"me"}},
	}

	require.Equal(t, 1, len(actual))
	require.EqualValues(t, expected, actual)
}

func TestSecredTestSuite(t *testing.T) {
	suite.Run(t, new(SecredTestSuite))
}
