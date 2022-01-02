package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type gameState struct {
	currentLocation Location
}

func setupDB() (*sql.DB, error) {
	os.Remove("data.db")

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadFile("data.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := string(buf)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	if err := CreateLocations(db); err != nil {
		panic(err)
	}

	if err := createItems(db); err != nil {
		panic(err)
	}

	return db, nil

}

func readInput(state *gameState, db *sql.DB, i string) {
	x := strings.Split(i, " ")
	c := commands[x[0]]
	switch c.(type) {
	case (func(*gameState, *sql.DB, string)):
		c.(func(*gameState, *sql.DB, string))(state, db, strings.Join(x[1:], " "))
	case func(*gameState):
		c.(func(*gameState))(state)
	}
}

func main() {

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	var state gameState

	state.currentLocation, err = getLocation("Living Room", db)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(state.currentLocation.Description)
	fmt.Printf("\n> ")

	readInput(&state, db, "move north")
	readInput(&state, db, "move up")
	readInput(&state, db, "take butter")
	readInput(&state, db, "push step ladder")
	readInput(&state, db, "move up")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		readInput(&state, db, scanner.Text())
		fmt.Printf("\n> ")
	}

}
