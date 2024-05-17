package lists_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xsadia/secred/pkg/utils/lists"
)

type S struct {
	id   string
	name string
	age  int
}

func TestMap(t *testing.T) {
	t.Run("Struct list to string list", func(t *testing.T) {
		structs := []S{
			S{id: "1", name: "user1", age: 22},
			S{id: "2", name: "user2", age: 22},
			S{id: "3", name: "user3", age: 22},
		}

		expected := []string{"1", "2", "3"}

		got := lists.Map(structs, func(elem S, _ int, _ []S) string {
			return elem.id
		})

		require.EqualValues(t, expected, got)
	})

	t.Run("Double", func(t *testing.T) {
		nums := []int{1, 2, 3}

		expected := []int{2, 4, 6}

		got := lists.Map(nums, func(elem int, _ int, _ []int) int {
			return elem * 2
		})

		require.EqualValues(t, expected, got)
	})
}

func TestReduce(t *testing.T) {
	t.Run("Sum all numbers", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		expected := 15

		got := lists.Reduce(nums, 0, func(acc int, elem int) int {
			return acc + elem
		})

		require.Equal(t, expected, got)
	})

	t.Run("Reduce into map", func(t *testing.T) {
		structs := []S{
			S{id: "1", name: "user1", age: 22},
			S{id: "2", name: "user2", age: 23},
			S{id: "3", name: "user3", age: 24},
		}

		expected := map[string]int{
			"user1": 22,
			"user2": 23,
			"user3": 24,
		}

		got := lists.Reduce(structs, map[string]int{}, func(acc map[string]int, elem S) map[string]int {
			acc[elem.name] = elem.age
			return acc
		})

		require.EqualValues(t, expected, got)
	})
}
