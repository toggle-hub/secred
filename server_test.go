package server_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/xsadia/secred/graph"
	"github.com/xsadia/secred/pkg/database"
)

func TestSecred(t *testing.T) {
	db, err := database.New("localhost", "root", "root", "secred_test", "disable")
	if err != nil {
		t.Errorf("failed to start test database %v", err.Error())
	}

	if err := database.Migrate(db); err != nil {
		t.Errorf("failed to migrate test database %v", err.Error())
	}

	c := client.New(handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}})))

	t.Run("createUser mutation", func(t *testing.T) {
		var resp map[string]interface{}
		c.MustPost(`mutation {
			createUser(input:{email:"fizi@gmail.com", name:"fizi", password:"123123"}) {
				name
				email
				__typename
			}
		}`, &resp)

		expected := map[string]interface{}{
			"name":       "fizi",
			"email":      "fizi@gmail.com",
			"__typename": "User",
		}
		actual := resp["createUser"].(map[string]interface{})
		require.Equal(t, expected["name"], actual["name"])
		require.Equal(t, expected["email"], actual["email"])
		require.Equal(t, expected["__typename"], actual["__typename"])
	})
}
