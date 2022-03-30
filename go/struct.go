package main

import "html/template"

type Templ struct { //Struct sent to api
	Artiste      []Artist
	Random       int
	Artistsearch []ArtistSearch
}
type Artist struct { //Struct used to get each artist's info
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Membersstr   string
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"` //1
	Location     struct { //1
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"` //2
		DatesLoc  struct { //2
			Dates []string `json:"dates"`
		}
	}
	JsString template.HTML
}

type ArtistSearch struct { //Struct used to get each artist's info
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Membersstr   string
	CreationDate int    `json:"creationDate"`
	FirstAlbum   string `json:"firstAlbum"`
}

type LattitudeLongitude struct { //Struct used to get map data for each artist
	Data []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"data"`
}
