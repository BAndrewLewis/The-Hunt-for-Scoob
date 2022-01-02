package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type item struct {
	ItemName    string   `json:"itemName"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
}

func getItem(itemName string, db *sql.DB) (item, error) {
	stmt, err := db.Prepare("select * from items where itemName = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	i := item{}
	var a string
	err = stmt.QueryRow(itemName).Scan(&i.ItemName, &i.Description, &a)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal([]byte(a), &i.Actions); err != nil {
		return item{}, err
	}

	return i, nil
}

func createItems(db *sql.DB) error {
	items := []item{

		{
			ItemName:    "pillows",
			Description: "Multi colored pillows are stacked on the couch cushions where you had been playing with your mom and Scoob earlier.",
			Actions:     []string{"move"},
		},
		{
			ItemName:    "couches",
			Description: "The couches are stacked with pillows. The seat cushions have a nice lip you could possibly climb up with",
			Actions:     []string{},
		},
		{
			ItemName:    "toys",
			Description: "Brightly painted blocks, a denim whale, and small toy piano are spread around.",
			Actions:     []string{"take"},
		},
		{
			ItemName:    "play chest",
			Description: "The play chest sits open in the corner. Upon further inspection you see that someone needs to clean off the jam and crumbs that, instead of eating, you smeared on it, but no Scoob.",
			Actions:     []string{},
		},
		{
			ItemName:    "scissors",
			Description: "Ooh, sharp!",
			Actions:     []string{"take"},
		},
		{
			ItemName:    "laptop",
			Description: "Your father is very focused on the bright light and pounding board of clackity-ness.",
			Actions:     []string{"take"},
		},
		{
			ItemName:    "breakable things",
			Description: "Ooh, breakable!",
			Actions:     []string{"take"},
		},
		{
			ItemName:    "butter",
			Description: "It looks very soft... and squishy...",
			Actions:     []string{"take"},
		},
		{
			ItemName:    "chairs",
			Description: "All of the chairs except for the one your dad is in are tucked under the table. One of the chairs has been pulled out just far enough that a small, agile, and nimble baby might be able to scramble up it. Your Dad is still in his chair.",
			Actions:     []string{"pull", "push"},
		},
		{
			ItemName:    "father",
			Description: "He’s sitting at the table, working. He seems very focused. He probably didn’t even notice you come in.",
			Actions:     []string{"talk"},
		},
		{
			ItemName:    "step ladder",
			Description: "The step ladder sits in the middle of the kitchen",
			Actions:     []string{"push", "pull"},
		},
	}
	for _, l := range items {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into items(itemName, description, actions) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		actions, err := json.Marshal(l.Actions)
		if err != nil {
			panic(err)
		}

		defer stmt.Close()
		_, err = stmt.Exec(l.ItemName, l.Description, actions)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}
	return nil
}

var doTakeFuncs = map[string]func(*gameState, *sql.DB, string){
	"butter": doTakeButter,
}

func doTakeButter(state *gameState, db *sql.DB, itemName string) {
	kitchen, err := getLocation("kitchen", db)
	if err != nil {
		panic(err)
	}
	state.currentLocation = kitchen
	fmt.Println()
	fmt.Println("You grab two big fistfulls of butter which squelches out between your fingers... Ah! Ingrid, haha!, good exploring, but butter is not something we play with. Let’s, ah! Don’t touch the computer! Let's wash your hands. Your father picks you up and carries you into the kitchen, stepping over the baby gate. Using warm water and soap your father washes the soap from your hands. Then places you on the Kitchen floor so he can attempt to wash off the butter you smeared all over his face and anywhere else you could reach. The counter forms a U around the room, closed off by the baby gate. You are still too short to see anything on the counter, but you can see a step ladder set up in the middle of the room.")
}

var doPushFuncs = map[string]func(*gameState, *sql.DB, string){
	"step ladder": doPushStepLadder,
}

func doPushStepLadder(state *gameState, db *sql.DB, itemName string) {
	newKitchen := Location{
		LocationName: "kitchen",
		Description:  "The counter forms a U around the room, closed off by the baby gate. You are still too short to see anything on the counter, the step ladder is propped against the cupboards.",
		Actions:      []string{},
		Items:        []string{"step ladder"},
		LinkedLocations: LinkedLocations{
			Up: "counter top",
		},
	}
	updateLocation(newKitchen, db)

	father := item{
		ItemName:    "father",
		Description: "He’s sitting at the table, working. He seems very focused. He probably didn’t even notice you come in.",
		Actions:     []string{"talk"},
	}

	updateItem(father, db)

	state.currentLocation = newKitchen

	fmt.Println()
	fmt.Println("You maneuver the ladder up against the cupboards. The countertop is now within your reach")
}

func updateItem(newItem item, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`
		UPDATE items
		SET description = ?,
			actions = ?
		WHERE itemName = ?`)
	if err != nil {
		log.Fatal(err)
	}

	actions, err := json.Marshal(newItem.Actions)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(newItem.Description, actions, newItem.ItemName)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
