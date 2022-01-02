package main

import (
	"database/sql"
	"fmt"
	"strings"
)

var commands = map[string]interface{}{
	"move":    doMove,
	"look":    doLook,
	"examine": doExamine,
	"take":    doTake,
	"push":    doPush,
}

func doMove(state *gameState, db *sql.DB, direction string) {
	direction = strings.ToLower(direction)
	var newLocation Location
	var err error
	switch direction {
	case "north":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.North, db)
		if err != nil {
			panic(err)
		}

	case "east":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.East, db)
		if err != nil {
			panic(err)
		}

	case "south":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.South, db)
		if err != nil {
			panic(err)
		}
	case "west":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.West, db)
		if err != nil {
			panic(err)
		}
	case "up":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.Up, db)
		if err != nil {
			panic(err)
		}
	case "down":
		newLocation, err = getLocation(state.currentLocation.LinkedLocations.Down, db)
		if err != nil {
			panic(err)
		}
	}
	state.currentLocation = newLocation
	fmt.Println()
	fmt.Println(state.currentLocation.Description)
}

func doLook(state *gameState) {
	fmt.Println()
	fmt.Println(state.currentLocation.Description)
}

func doExamine(state *gameState, db *sql.DB, itemName string) {
	// TODO update this so that it checks the players inventory too
	if v := contains(strings.ToLower(itemName), state.currentLocation.Items); v {
		fmt.Println("That item doesn't seem to be in the room")
		return
	}

	item, err := getItem(itemName, db)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(item.Description)
}

func doTake(state *gameState, db *sql.DB, itemName string) {
	// make sure the item is in the room
	if v := contains(strings.ToLower(itemName), state.currentLocation.Items); !v {
		fmt.Println("Can't find that")
		return
	}

	// make sure the item can be pushed
	item, err := getItem(itemName, db)
	if err != nil {
		panic(err)
	}
	if v := contains("take", item.Actions); !v {
		fmt.Println("Can't take that")
		return
	}

	// call the items take function
	doTakeFuncs[itemName](state, db, itemName)

}

func doPush(state *gameState, db *sql.DB, itemName string) {
	// make sure the item is in the room
	if v := contains(strings.ToLower(itemName), state.currentLocation.Items); !v {
		fmt.Println("Can't find that")
		return
	}
	// make sure the item can be pushed
	item, err := getItem(itemName, db)
	if err != nil {
		panic(err)
	}
	if v := contains("push", item.Actions); !v {
		fmt.Println("Can't push that")
		return
	}

	// call the items push function
	doPushFuncs[itemName](state, db, itemName)
}

func doDrop() {}

func contains(item string, items []string) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}
	return false
}
