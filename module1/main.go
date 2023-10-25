package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type CS2Skin struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Wear        string `json:"wear"`
	Pattern     string `json:"pattern"`
	Side        string `json:"side"`
	WeaponType  string `json:"weapon_type"`
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу, если ее нет
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cs2_skins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		price INTEGER,
		wear TEXT,
		pattern TEXT,
		side TEXT,
		weapon_type TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var skin CS2Skin
		err := json.NewDecoder(r.Body).Decode(&skin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO cs2_skins (name, price, wear, pattern, side, weapon_type) VALUES (?, ?, ?, ?, ?, ?)",
			skin.Name, skin.Price, skin.Wear, skin.Pattern, skin.Side, skin.WeaponType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Skin %s added successfully!", skin.Name)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price, wear, pattern, side, weapon_type FROM cs2_skins")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var skins []CS2Skin
		for rows.Next() {
			var skin CS2Skin
			err := rows.Scan(&skin.ID, &skin.Name, &skin.Price, &skin.Wear, &skin.Pattern, &skin.Side, &skin.WeaponType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			skins = append(skins, skin)
		}

		json.NewEncoder(w).Encode(skins)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}