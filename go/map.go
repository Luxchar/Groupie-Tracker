package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

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
	s, _ := ioutil.ReadFile("../static/assets/js/mapjs.txt")
	script := string(s)
	k := `<script>
	var LocationsForMap = [
		` + strings.Join(slice, "\n") + script
	artist.JsString = template.HTML(k)
}
