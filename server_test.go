package server_test

import (
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
	"github.com/xsadia/secred/pkg/database"
)

type Me struct {
	Name     string
	Email    string
	Typename string `json:"__typename"`
}
type CreateUser struct {
	Me Me
}

type Response struct {
	CreateUser CreateUser
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

func (suite *SecredTestSuite) TestCreateUserSuccess() {
	t := suite.T()

	var actual Response
	suite.client.MustPost(`mutation {
			createUser(input:{email:"fizi@gmail.com", name:"fizi", password:"123123"}) {
				me {
					name
					email
					__typename
				}
			}
		}`, &actual)

	expected := Response{
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
	suite.db.Exec("INSERT INTO users (name, email, password) values ($1, $2, $3)", "fizi", "fizi@gmail.com", "123123")
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

func TestSecredTestSuite(t *testing.T) {
	suite.Run(t, new(SecredTestSuite))
}
