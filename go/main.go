package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var templatesDir = os.Getenv("TEMPLATES_DIR")
var artist Artist //creation instance struct artist
var art Templ     //template to send to front

func tracker(w http.ResponseWriter, r *http.Request) { //function that starts when tracker or artist page is loaded
	var (
		url         string //url for making requests
		Randparam   = r.URL.Query()["RandomArtist"]
		apparitionP = r.URL.Query()["apparition"]
		albumP      = r.URL.Query()["album"]
		membersP    = r.URL.Query()["members"]
		locationsP  = r.URL.Query()["location"]
	)
	var emptyart Templ //reset struct to not have the same artist piling up
	art = emptyart

	rand.Seed(time.Now().UnixNano()) //random number to pick random artist
	art.Random = 1 + rand.Intn(51-0)

	var generated []string //stock random integer for randomizer tracker page
	url = "https://groupietrackers.herokuapp.com/api/artists/"
	//-------------------------- FETCH RESSOURCES API JSON
	if Randparam != nil { //user wants random artist page
		url += Randparam[0]
		fetchartist(url, apparitionP, albumP, membersP, locationsP, false)
	} else { // user wants page tracker with every artist on it
		for i := 1; i < 11; i++ {
			rand.Seed(time.Now().UnixNano()) //random number to pick random artist
			random := strconv.Itoa(1 + rand.Intn(51-0))
			for stringInSlice(random, generated) { //while number is not unique, picks another one
				random = strconv.Itoa(1 + rand.Intn(51-0))
			}
			generated = append(generated, random)
			url += random
			fetchartist(url, apparitionP, albumP, membersP, locationsP, false)
			url = "https://groupietrackers.herokuapp.com/api/artists/"
		}
		var artistsearch ArtistSearch
		for i := 1; i < 53; i++ { //for searchbar research
			rand.Seed(time.Now().UnixNano()) //random number to pick random artist
			url += strconv.Itoa(i)
			err := json.Unmarshal([]byte(request(url)), &artistsearch)
			if err != nil {
				fmt.Print("error when encoding the struct")
			}
			url = "https://groupietrackers.herokuapp.com/api/artists/"
			art.Artistsearch = append(art.Artistsearch, artistsearch)
		}
	}
	(template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/tracker.html")))).Execute(w, art)
}

func research(apparitionP, albumP, membersP, locationsP []string) bool { //does the user want to research
	if apparitionP != nil || albumP != nil || membersP != nil || locationsP != nil {
		return true
	}
	return false
}

func criteria(apparitionP, albumP, membersP, locationsP []string) {
	if research(apparitionP, albumP, membersP, locationsP) { //if user has chosen some criteria
		if search(artist, apparitionP, albumP, membersP, locationsP) { //check whether artist in api is eligible to show
			art.Artiste = append(art.Artiste, artist) //put it on the page
		}
	} else { //no criteria, just display everything
		art.Artiste = append(art.Artiste, artist)
	}
}

func fetchartist(url string, apparitionP, albumP, membersP, locationsP []string, needmap bool) { //fetch an artist from the api
	err, _, _ := json.Unmarshal([]byte(request(url)), &artist),
		json.Unmarshal([]byte(request(artist.Locations)), &artist.Location),
		json.Unmarshal([]byte(request(artist.Location.Dates)), &artist.Location.DatesLoc)
	str := ""
	for _, v := range artist.Members { //prints members without []
		str += v + " "
	}
	artist.Membersstr = str
	if err != nil {
		fmt.Print("error when encoding the struct")
	}
	if needmap { //if user tries to print an artist page he needs the map
		go JscriptStr()
	}
	criteria(apparitionP, albumP, membersP, locationsP)
}

func search(artist Artist, apparitionP []string, albumP []string, membersP []string, locationsP []string) bool { //handles the comparisons with the users criteria to the artists
	if apparitionP != nil {
		if strconv.Itoa(artist.CreationDate) != apparitionP[0] && apparitionP[0] != "1945" {
			return false
		}
	}
	if albumP != nil {
		if artist.FirstAlbum != albumP[0] && albumP[0] != "1945" {
			return false
		}
	}
	if membersP != nil {
		if strconv.Itoa(len(artist.Members)) != membersP[0] && membersP[0] != "" {
			return false
		}
	}
	if locationsP != nil {
		if !strings.Contains(artist.Location.Locations[0], locationsP[0]) {
			return false
		}
	}
	return true
}

func request(url string) []byte { //simple get request to get the api data
	req, _ := http.Get(url)
	body, _ := ioutil.ReadAll(req.Body)
	return body
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func artistt(w http.ResponseWriter, r *http.Request) { //user wants specific artist page
	var emptyart Templ //reset struct to not have the same artist piling up
	art = emptyart
	param := r.URL.Query()["artist"]
	url := "https://groupietrackers.herokuapp.com/api/artists/"
	url += param[0]
	var none []string
	fetchartist(url, none, none, none, none, true)

	(template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/artist.html")))).Execute(w, art)
}

func main() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.HandleFunc("/pages/tracker", tracker)
	http.HandleFunc("/pages/artist", artistt)
	fmt.Printf("Started server successfully on http://localhost:8089/\n")
	http.ListenAndServe(":8089", nil)
}
