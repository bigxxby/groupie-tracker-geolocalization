package bigxxby

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	Id            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	DateLocations map[string][]string
}

type Index struct {
	Index []Relation `json:"index"`
}

type Relation struct {
	DateLocations map[string][]string `json:"datesLocations"`
}

func GetContent() []Artist {
	index := getRelations().Index
	artists := getArtists()
	for i := 0; i < len(artists); i++ {
		artists[i].DateLocations = index[i].DateLocations
	}
	return artists
}

func getArtists() []Artist {
	content, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	files, err := ioutil.ReadAll(content.Body)
	if err != nil {
		return nil
	}
	defer content.Body.Close()
	var artists []Artist
	err = json.Unmarshal(files, &artists)
	if err != nil {
		log.Println("error", err.Error())
		return nil
	}

	return artists
}

func getRelations() Index {
	var index Index
	content, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Println("Get error: ", err.Error())
		return index
	}
	files, err := ioutil.ReadAll(content.Body)
	if err != nil {
		return index
	}
	defer content.Body.Close()
	err = json.Unmarshal(files, &index)
	if err != nil {
		log.Println("error", err.Error())
		return index
	}
	return index
}
