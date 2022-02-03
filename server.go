package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type article struct {
	Author      string
	ImageUrl    string
	Content     string
	Date        string
	ReadMoreUrl string `json:"readMoreUrl"`
	Title       string
	Time        string
}

type articlesList struct {
	Category string
	Data     []article
}

type returnData struct {
	Articles []article
	Title    string
}

type mockData struct {
	Data []article
}

func getArticles() ([]article, error) {
	response, _ := http.Get("https://inshortsapi.vercel.app/news?category=technology")
	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var list articlesList
	err := json.Unmarshal(bytes, &list)

	if err != nil {
		return nil, err
	}

	return list.Data, nil
}

func getMockArticles() ([]article, error) {
	jsonFile, err := os.Open("./mock-data/mock.json")

	if err != nil {
		return nil, err
	}

	mockBytes, _ := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	var mock mockData
	json.Unmarshal(mockBytes, &mock)

	return mock.Data, nil
}

func articlesHandler(response http.ResponseWriter, request *http.Request) {
	// articles, err := getMockArticles()
	articles, err := getArticles()

	if err != nil {
		panic(err)
	}

	view, _ := template.ParseFiles("./views/articles.html")
	data := returnData{Articles: articles, Title: "Technology InShorts"}

	view.Execute(response, data)
}

func main() {
	styles := http.FileServer(http.Dir("./views/stylesheets"))
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	http.HandleFunc("/articles", articlesHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
