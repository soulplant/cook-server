package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// A User that can log in.
type User struct {
	ID    int64
	Name  string
	Email string
}

// Instruction is a step in a recipe.
type Instruction struct {
	ID       int64
	RecipeID int64
	Step     int
	Text     string
}

// Aisle is a supermarket aisle where Ingredients can be found.
type Aisle struct {
	ID   int64
	Name string
}

// Ingredient definition.
type Ingredient struct {
	ID   int64
	Name string

	// The aisle where this ingredient can be found.
	Aisle   Aisle
	AisleID int64
}

// Unit of measurement.
type Unit struct {
	ID   int64
	Name string
	// True if a measurement, e.g. g or kg which should be rendered immediately
	// after the count.
	Measurement bool
}

// RecipeIngredient is an ingredient as part of a recipe. This is as opposed to
// Ingredient, which is the ingredient itself, this type includes information
// about the quantity and preparation of the ingredient in this particular
// recipe.
type RecipeIngredient struct {
	Recipe       Recipe
	RecipeID     int64
	Ingredient   Ingredient
	IngredientID int64
	Unit         Unit
	UnitID       int64

	// How many units of the ingredient, measured in Units.
	Quantity int

	// Instructions for how to prepare this ingredient.
	Preparation string
}

// Recipe represents a single entry in a cookbook.
type Recipe struct {
	ID           int64
	Name         string
	Ingredients  []RecipeIngredient
	Instructions []Instruction
	Author       User
	AuthorID     int64
}

var allTables = []interface{}{
	&User{},
	&Instruction{},
	&Recipe{},
	&Ingredient{},
	&RecipeIngredient{},
	&Aisle{},
	&Unit{},
}

func addUser(b *gorm.DB, name string) *User {
	u := User{
		Name: name,
	}
	if err := b.Save(&u).Error; err != nil {
		log.Fatal("Failed to save user", name, err)
	}
	return &u
}

func addRecipe(db *gorm.DB, userID int64, recipeName string, steps []string) *Recipe {
	var ins []Instruction
	for i, s := range steps {
		ins = append(ins, Instruction{Step: i, Text: s})
	}
	r := Recipe{
		Name:         recipeName,
		AuthorID:     userID,
		Instructions: ins,
	}
	db.Save(&r)
	return &r
}

func initTestData(db *gorm.DB) {
	if err := db.DropTableIfExists(allTables...).Error; err != nil {
		log.Fatal("Failed to drop tables", err)
	}
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		log.Fatal("Failed to automigrate", err)
	}
	j := addUser(db, "james")
	s := addUser(db, "steve")
	addRecipe(db, j.ID, "Comfort Pasta", []string{"step1", "step2"})
	r := addRecipe(db, s.ID, "Simple Curry", []string{"step1", "step 2", "step3"})
	getRecipe(db, r.ID)
}

func getRecipe(db *gorm.DB, recipeID int64) (*Recipe, error) {
	r := Recipe{ID: recipeID}
	err := db.Preload("Instructions").Preload("Ingredients.Ingredient").Find(&r).Error
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func findRecipes(db *gorm.DB) ([]*Recipe, error) {
	r := []*Recipe{}
	err := db.Preload("Ingredients.Ingredient").Preload("Instructions").Order("id").Find(&r).Error
	return r, err
}

func main2() {
	db, err := sql.Open("postgres", "sslmode=disable")
	if err != nil {
		log.Fatal("couldn't open db", err)
	}
	_, _ = db.Exec("drop database test2")
	_, err = db.Exec("create database test2")
	if err != nil {
		log.Fatal("couldn't create test2", err)
	}
	fmt.Println("created test2")
	err = db.Close()
	if err != nil {
		log.Fatal("Couldn't close db", err)
	}
	gdb, err := gorm.Open("postgres", "database=test2 sslmode=disable")
	if err != nil {
		log.Fatal("Couldn't gorm test2", err)
	}
	gdb.AutoMigrate(allTables...)
	u := User{
		Name: "james",
	}
	if err := gdb.Save(&u).Error; err != nil {
		panic(err)
	}
	if err := gdb.Model(&u).Scan(&u).Error; err != nil {
		log.Fatal("Couldn't scan", err)
	}
	fmt.Printf("Got user %v\n", u)
}

func main() {
	fmt.Printf("drivers %v\n", sql.Drivers())
	db, err := gorm.Open("postgres", "database=test sslmode=disable")
	if err != nil {
		panic(err)
	}
	initTestData(db)
	var b User
	if err := db.Model(&User{}).Where("name = ?", "james").Scan(&b).Error; err != nil {
		log.Fatal("Failed to query users", err)
	}

	fmt.Printf("HI %v\n", b)
}
