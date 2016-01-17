package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	bfs "github.com/kbuzsaki/wikidegree/bfs"
	api "github.com/kbuzsaki/wikidegree/api"
)

func lookup(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	start := values.Get("start")
	end := values.Get("end")

	path, err := lookupPath(start, end)
	if err != nil {
		log.Print(err)
		io.WriteString(writer, "Error: " + err.Error())
	} else {
		pathBytes, _ := json.Marshal(&path)
		io.WriteString(writer, string(pathBytes))
	}
}

func lookupPath(start, end string) (api.TitlePath, error) {
	// valiate start and end titles exist
	if start == "" || end == "" {
		return nil, errors.New("start and end parameters required")
	}
	start = api.EncodeTitle(start)
	end = api.EncodeTitle(end)

	// load the page loader, currently only bolt
	pageLoader, err := api.GetBoltPageLoader()
	if err != nil {
		return nil, err
	}
	defer pageLoader.Close()

	// validate that the start page exists and has links
	startPage, err := pageLoader.LoadPage(start)
	if err != nil {
		return nil, err
	}
	if len(startPage.Links) == 0 {
		return nil, errors.New("start page has no links!")
	}

	// validate that the end page exists
	_, err = pageLoader.LoadPage(end)
	if err != nil {
		return nil, err
	}

	// actually find the path using bfs
	pathFinder := bfs.GetBfsPathFinder(pageLoader)
	return pathFinder.FindPath(start, end)
}

func main() {
	http.HandleFunc("/", lookup)
	http.ListenAndServe(":8000", nil)
}