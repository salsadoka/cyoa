package main

import (
	"flag"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/salsadoka/cyoa"
	"net/http"
	"fmt"
)

func main(){
	port := flag.Int("port", 3000, "the port to start the CYOA web app on")
	fileName := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	story := buildStoryArcsMap(*fileName)

	h := cyoa.NewHandler(story)
	log.Printf("Starting the server at: %d \n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

}


func buildStoryArcsMap( fileName string) cyoa.Story {
	file, err := os.Open(fileName)
	if nil != err {
		log.Fatal("Failed to open file with error: %s", err.Error())
	}
	j, err := ioutil.ReadAll(file)
	if nil != err {
		log.Fatal("Failed to read file with error: %s", err.Error())
	}
	storyArcs := make(cyoa.Story)
	err = json.Unmarshal(j, &storyArcs)
	if nil != err {
		log.Fatalf("Failed to unmarshall json with error: %s", err.Error())
	}
	return storyArcs
}