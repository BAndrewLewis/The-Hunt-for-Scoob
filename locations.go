package main

import (
	"database/sql"
	"encoding/json"
	"log"
)

type Location struct {
	LocationName    string          `json:"locationName"`
	Description     string          `json:"description"`
	Actions         []string        `json:"actions"`
	Items           []string        `json:"items"`
	LinkedLocations LinkedLocations `json:"linkedLocations"`
}

type LinkedLocations struct {
	North string `json:"north,omitempty"`
	East  string `json:"east,omitempty"`
	South string `json:"south,omitempty"`
	West  string `json:"west,omitempty"`
	Up    string `json:"up,omitempty"`
	Down  string `json:"down,omitempty"`
}

func getLocation(locationName string, db *sql.DB) (Location, error) {
	stmt, err := db.Prepare("select * from locations where locationName = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	l := Location{}
	var a string
	var i string
	var ll string
	err = stmt.QueryRow(locationName).Scan(&l.LocationName, &l.Description, &a, &i, &ll)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal([]byte(a), &l.Actions); err != nil {
		return Location{}, err
	}

	if err := json.Unmarshal([]byte(i), &l.Items); err != nil {
		return Location{}, err
	}

	if err := json.Unmarshal([]byte(ll), &l.LinkedLocations); err != nil {
		return Location{}, err
	}

	return l, nil
}

func CreateLocations(db *sql.DB) error {
	locations := []Location{

		{
			LocationName: "Living Room",
			Description:  "Looking around the living room you can see several couches, stacked with pillows, and, various toys spilling from the play chest and across the floor. Where did you last have your Scoob? There are exits to the North and East",
			Actions:      []string{},
			Items:        []string{"couches", "pillows", "toys", "play chest"},
			LinkedLocations: LinkedLocations{
				North: "Dining Room",
				East:  "Baby Gate",
			},
		},

		{
			LocationName: "Dining Room",
			Description:  "The table in the dining room is surrounded by 4 chairs, all tucked in, except for one which is occupied by your father. You know there are things on the table, there are always things on the table, but you can’t see from down here.",
			Actions:      []string{},
			Items:        []string{"chairs"},
			LinkedLocations: LinkedLocations{
				Up:    "Table",
				South: "Living Room",
			},
		},

		{
			LocationName: "Baby Gate",
			Description:  "Hmmm. The gate is too tall for you to even be able to pull yourself up on, and is firmly attached to the wall bordering the entryway into the Kitchen. No amount of rattling, or screaming has ever gotten it to budge.",
			Actions:      []string{},
			Items:        []string{},
			LinkedLocations: LinkedLocations{
				West: "Living Room",
			},
		},
		{
			LocationName: "Table",
			Description:  "You can now see the spread of forbidden objects on the table top. Scissors, butter, Dad’s laptop, and an assortment of breakable things, but no Scoob.",
			Actions:      []string{},
			Items:        []string{"scissors", "butter", "laptop", "breakable things"},
			LinkedLocations: LinkedLocations{
				Down: "Dining Room",
			},
		},
		{
			LocationName:    "kitchen",
			Description:     "The counter forms a U around the room, closed off by the baby gate. You are still too short to see anything on the counter, but you can see a step ladder set up in the middle of the room.",
			Actions:         []string{},
			Items:           []string{"step ladder"},
			LinkedLocations: LinkedLocations{},
		},
		{
			LocationName: "counter top",
			Description:  "You crawl along the cold granite counter top on your hands and knees. Your father is still struggling to get the butter off his face at the sink. Various snacks, and kitchen appliances line the counter on the far side of the room. Unfortunately, the way over is blocked by the sink and your father. On your side of the counter is a pile of dishes.",
			Actions:      []string{},
			Items:        []string{"sink", "father", "kitchen appliances", "pile of dishes", "bowls", "snacks"},
			LinkedLocations: LinkedLocations{
				Down: "kitchen",
			},
		},
	}
	for _, l := range locations {

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into locations(locationName, description, actions, items, linkedLocations) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}

		actions, err := json.Marshal(l.Actions)
		if err != nil {
			panic(err)
		}

		items, err := json.Marshal(l.Items)
		if err != nil {
			panic(err)
		}

		linkedLocations, err := json.Marshal(l.LinkedLocations)
		if err != nil {
			panic(err)
		}

		defer stmt.Close()
		_, err = stmt.Exec(l.LocationName, l.Description, actions, items, linkedLocations)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}
	return nil
}

func updateLocation(newLocation Location, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`
		UPDATE locations
		SET description = ?,
			actions = ?,
			items = ?,
			linkedLocations = ? 
		WHERE locationName = ?`)
	if err != nil {
		log.Fatal(err)
	}

	actions, err := json.Marshal(newLocation.Actions)
	if err != nil {
		panic(err)
	}

	items, err := json.Marshal(newLocation.Items)
	if err != nil {
		panic(err)
	}

	linkedLocations, err := json.Marshal(newLocation.LinkedLocations)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(newLocation.Description, actions, items, linkedLocations, newLocation.LocationName)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
