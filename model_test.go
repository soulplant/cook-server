package main

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func newTestGorm(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open("postgres", "database=test2 sslmode=disable")
	if err != nil {
		t.Fatal("Couldn't get test gorm", err)
	}
	gdb.DropTableIfExists(allTables...)
	if err := gdb.AutoMigrate(allTables...).Error; err != nil {
		t.Fatal("Failed automigrate", err)
	}
	return gdb
}

func TestStuff(t *testing.T) {
	db := newTestGorm(t)
	j := User{Name: "james"}
	db.Save(&j)
	r := Recipe{Author: j, Name: "Comfort Pasta", Instructions: []Instruction{
		{Step: 1, Text: "Do something"},
		{Step: 2, Text: "Do something else"},
		{Step: 3, Text: "Do something else again!"},
	},
		Ingredients: []RecipeIngredient{
			{
				Quantity: 1,
				Unit:     Unit{Name: "cup", Measurement: false},
				Ingredient: Ingredient{
					Name:  "rice",
					Aisle: Aisle{Name: "Rice Aisle"},
				},
				Preparation: "steamed",
			},
		},
	}
	db.Save(&r)
	nr, err := getRecipe(db, r.ID)
	if err != nil {
		t.Fatal("Failed to get recipe", err)
	}
	if len(nr.Instructions) != 3 {
		t.Errorf("Expected 3 instructions, got %v", nr.Instructions)
	}
	if len(nr.Ingredients) != 1 {
		t.Fatalf("Expected 1 ingredient, got %v", nr.Ingredients)
	}
	if nr.Ingredients[0].Ingredient.Name != "rice" {
		t.Fatalf("Expected ingredient name to be 'rice', but was %v", nr.Ingredients[0].Ingredient.Name)
	}

	r2 := Recipe{Author: j, Name: "Salad sandwich", Instructions: []Instruction{
		{Step: 1, Text: "Hey"},
	}}
	db.Save(&r2)
	nr2, err := getRecipe(db, r2.ID)
	if err != nil {
		t.Fatal("Failed to get recipe", err)
	}
	if len(nr2.Instructions) != 1 {
		t.Errorf("Expected 1 instruction, got %v", nr2.Instructions)
	}

	rs, err := findRecipes(db)
	if err != nil {
		t.Fatal("Failed to find recipes", err)
	}
	if len(rs) != 2 {
		t.Fatalf("Expected 2 recipes, got %v", rs)
	}
	if rs[0].Ingredients[0].Ingredient.Name != "rice" {
		t.Errorf("Expected 'rice', got %v", rs[0].Ingredients[0].Ingredient)
	}
}
