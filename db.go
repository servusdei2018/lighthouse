package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

type Listing struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
	Url string `json:"url"`
}

func (l *Listing) String() string {
	return fmt.Sprintf("Name: %s, Host: %s, Port: %v, Website: %s", l.Name, l.Host, l.Port, l.Url)
}

// DB implements a MUD's database.
type DB struct {
	Listings map[int]Listing `json:"listings"`
	ListingKeys []int
}

// Export saves the MUD database to disk.
func (d *DB) Export(path string) (err error) {
	f, err := os.Create(path)
	if err != nil { return }
	defer f.Close()

	// Delete prior data.
	f.Truncate(0)

	// Marshal DB to JSON.
	data, err := json.Marshal(d)
	if err != nil { return }

	// Output JSON to file.
	_, err = f.Write(data)
	return
}

// Import loads MUD database from disk.
func (d *DB) Import(path string) (err error) {
	d.Listings = make(map[int]Listing)

	// Load JSON from file.
	data, err := ioutil.ReadFile(path)
	if err != nil { return }

	// Parse JSON and populate DB.
	err = json.Unmarshal(data, d)
	if err != nil { return }

	// Generate slice of keys.
	for key, _ := range d.Listings {
		d.ListingKeys = append(d.ListingKeys, key)
	}
	return
}

// RandomListing returns a random MUD listing from the database.
func (d *DB) RandomListing() Listing {
	return d.Listings[rand.Intn(len(d.ListingKeys))]
}
