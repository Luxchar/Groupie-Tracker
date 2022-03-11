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

//Js dans un fichier
type Templ struct { //Struct sent to api
	Artiste []Artist
	Random  int
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

type LattitudeLongitude struct { //Struct used to get map data for each artist
	Data []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"data"`
}

var templatesDir = os.Getenv("TEMPLATES_DIR")
var artist Artist //creation instance struct artist

func tracker(w http.ResponseWriter, r *http.Request) { //function that starts when tracker or artist page is loaded
	var (
		art         Templ
		url         string //url for making requests
		selectart   bool
		param       = r.URL.Query()["artist"]
		Randparam   = r.URL.Query()["RandomArtist"]
		apparitionP = r.URL.Query()["apparition"]
		albumP      = r.URL.Query()["album"]
		membersP    = r.URL.Query()["members"]
		locationsP  = r.URL.Query()["location"]
	)
	rand.Seed(time.Now().UnixNano()) //random number to pick random artist
	art.Random = 1 + rand.Intn(51-0)

	var research bool
	var iteration int
	var generated []string //stock random integer for randomizer tracker page
	if apparitionP != nil || albumP != nil || membersP != nil || locationsP != nil {
		research = true
		iteration = 53
	} else {
		iteration = 11
	}

	for i := 1; i < iteration; i++ {
		url = "https://groupietrackers.herokuapp.com/api/artists/"
		artist.Members = nil
		//-------------------------- FETCH RESSOURCES API JSON
		if param != nil { //user wants specific artist page
			url += param[0]
			selectart = true
		} else if Randparam != nil { //user wants random artist page
			url += Randparam[0]
			selectart = true
		} else { //user wants page tracker with every artist on it
			rand.Seed(time.Now().UnixNano()) //random number to pick random artist
			random := strconv.Itoa(1 + rand.Intn(51-0))
			for stringInSlice(random, generated) {
				random = strconv.Itoa(1 + rand.Intn(51-0))
			}
			generated = append(generated, random)
			url += random
			selectart = false
		}
		s, _, _ := json.Unmarshal([]byte(request(url)), &artist),
			json.Unmarshal([]byte(request(artist.Locations)), &artist.Location),
			json.Unmarshal([]byte(request(artist.Location.Dates)), &artist.Location.DatesLoc)
		str := ""
		for _, v := range artist.Members { //prints members without []
			str += v + " "
		}
		artist.Membersstr = str
		if s != nil {
			fmt.Print("error when encoding the struct")
		}
		//-------------------------- append a la struct envoyée a l'api l'instance en cours
		if research { //if user has chosen some criteria
			if search(artist, apparitionP, albumP, membersP, locationsP) { //check whether artist in api is eligible to show
				art.Artiste = append(art.Artiste, artist) //put it on the page
			}
		} else { //no criteria, just display everything
			art.Artiste = append(art.Artiste, artist)
		}

		if selectart { //if artist is selected print the page of the artist and stop
			JscriptStr()
			(template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/artist.html")))).Execute(w, art)
			return
		}
		url = ""
	}
	(template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/tracker.html")))).Execute(w, art)
}

func search(artist Artist, apparitionP []string, albumP []string, membersP []string, locationsP []string) bool { //handles the comparisons with the users criteria to the artists
	if apparitionP != nil {
		if strconv.Itoa(artist.CreationDate) == apparitionP[0] && apparitionP[0] != "1945" {
			return true
		}
	}
	if albumP != nil {
		if artist.FirstAlbum == albumP[0] && albumP[0] != "1945" {
			return true
		}
	}
	if membersP != nil {
		if strconv.Itoa(len(artist.Members)) == membersP[0] && membersP[0] != "" {
			return true
		}
	}
	if locationsP != nil {
		if strings.Contains(artist.Location.Locations[0], locationsP[0]) {
			return true
		}
	}
	return false
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

func JscriptStr() { //handles the map on the artist page
	var loc LattitudeLongitude
	var slice []string
	str := ""
	for limite, i := range artist.Location.Locations { //Retire cas particuliers
		for _, j := range i {
			if j == '_' {
				j = ' '
			}
			str += string(j)
		}
		json.Unmarshal([]byte(request("http://api.positionstack.com/v1/forward?access_key=9a6d5681ba2b143da463543ee17cf96e&query="+str+"&limit=1")), &loc)
		if len(artist.Location.Locations)-1 == limite { //créations jsstring
			slice = append(slice, "['"+str+"',"+fmt.Sprintf("%f", loc.Data[0].Latitude)+", "+fmt.Sprintf("%f", loc.Data[0].Longitude)+","+fmt.Sprintf("%q", artist.Location.DatesLoc.Dates[limite][1:])+"],")
		} else {
			slice = append(slice, `['`+str+"',"+fmt.Sprintf("%f", loc.Data[0].Latitude)+", "+fmt.Sprintf("%f", loc.Data[0].Longitude)+","+fmt.Sprintf("%q", artist.Location.DatesLoc.Dates[limite][1:])+"],")
		}
		str = ""
	}
	k := `<script>
	var LocationsForMap = [
		` + strings.Join(slice, "\n") + `
		];

	// Initialize and add the map
	function initMap() {
	  // The map, centered 
	  const map = new google.maps.Map(document.getElementById("map"), {
		zoom: 1,
		center: new google.maps.LatLng(36.91, 1.64),
	  });

	  var infowindow = new google.maps.InfoWindow();

	  for (i = 0; i < LocationsForMap.length; i++) {
		const marker = new google.maps.Marker({
		  position: new google.maps.LatLng(LocationsForMap[i][1], LocationsForMap[i][2]),
		  map: map,
		});
  
		google.maps.event.addListener(marker, 'click', (function(marker, i) {
		  return function() {
			infowindow.setContent(LocationsForMap[i][0]+" "+LocationsForMap[i][3]);
			infowindow.open(map, marker);
		  }
		})(marker, i));
	  }
	}
  </script>`
	artist.JsString = template.HTML(k)
}

func main() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.HandleFunc("/pages/tracker", tracker)
	http.HandleFunc("/pages/artist", tracker)
	fmt.Printf("Started server successfully on http://localhost:8089/\n")
	http.ListenAndServe(":8089", nil)
}
