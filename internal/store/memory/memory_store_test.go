package memory

import (
	"brewday/internal/recipe"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateID(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name       string
		RecipeName string
		Expected   string
	}
	testCases := []testCase{
		{
			RecipeName: "TestRecipe1",
			Expected:   "5465737452656369706531", // Hex encoding of "TestRecipe1"
		},
		{
			RecipeName: "AnotherRecipe",
			Expected:   "416e6f74686572526563697065", // Hex encoding of "AnotherRecipe"
		},
		{
			RecipeName: "12345",
			Expected:   "3132333435", // Hex encoding of "12345"
		},
	}
	store := NewMemoryStore()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(tc.Expected, store.CreateID(tc.RecipeName))
		})
	}
}

func TestStore(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name    string
		Recipes []*recipe.Recipe
	}
	testCases := []testCase{
		{
			Name:    "Single recipe",
			Recipes: []*recipe.Recipe{{Name: "recipe1"}},
		},
		{
			Name:    "Multiple recipes",
			Recipes: []*recipe.Recipe{{Name: "recipe1"}, {Name: "recipe2"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			store := NewMemoryStore()
			m := make(map[string]*recipe.Recipe)
			for _, re := range tc.Recipes {
				id, _ := store.Store(re)
				m[id] = re
			}
			require.Equal(len(tc.Recipes), len(store.recipes))
			for id, re := range m {
				require.Equal(re, store.recipes[id])
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	require := require.New(t)
	store := NewMemoryStore()
	store.recipes["id1"] = &recipe.Recipe{Name: "recipe1"}
	type testCase struct {
		Name     string
		ID       string
		Expected *recipe.Recipe
		Error    bool
	}
	testCases := []testCase{
		{
			Name:     "Existing recipe",
			ID:       "id1",
			Expected: &recipe.Recipe{Name: "recipe1"},
			Error:    false,
		},
		{
			Name:     "Non-existing recipe",
			ID:       "id2",
			Expected: nil,
			Error:    true,
		},
		{
			Name:     "Empty ID",
			ID:       "",
			Expected: nil,
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			re, err := store.Retrieve(tc.ID)
			if tc.Error {
				require.Nil(re)
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Expected, re)
			}
		})
	}
}

func TestList(t *testing.T) {
	require := require.New(t)
	storeSingle := NewMemoryStore()
	storeSingle.recipes["id1"] = &recipe.Recipe{Name: "recipe1"}
	storeMultiple := NewMemoryStore()
	storeMultiple.recipes["id1"] = &recipe.Recipe{Name: "recipe1"}
	storeMultiple.recipes["id2"] = &recipe.Recipe{Name: "recipe2"}
	storeEmpty := NewMemoryStore()
	type testCase struct {
		Name     string
		Store    *MemoryStore
		Expected []*recipe.Recipe
	}
	testCases := []testCase{
		{
			Name:     "Single recipe",
			Store:    storeSingle,
			Expected: []*recipe.Recipe{{Name: "recipe1"}},
		},
		{
			Name:     "Multiple recipes",
			Store:    storeMultiple,
			Expected: []*recipe.Recipe{{Name: "recipe1"}, {Name: "recipe2"}},
		},
		{
			Name:     "Empty store",
			Store:    storeEmpty,
			Expected: []*recipe.Recipe{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			recipes, err := tc.Store.List()
			require.NoError(err)
			require.ElementsMatch(tc.Expected, recipes)
		})
	}
}
