package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/andrewslotin/doppelganger/github"
)

var (
	reposTemplate = template.Must(template.ParseFiles("templates/repos/index.html.template"))
)

type ReposHandler struct {
	repositories github.RepositoryService
}

func NewReposHandler(repositoryService github.RepositoryService) *ReposHandler {
	return &ReposHandler{
		repositories: repositoryService,
	}
}

func (handler *ReposHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	if repoName := req.FormValue("repo"); repoName != "" {
		NewRepoClient(handler.repositories).ServeHTTP(w, req)
		return
	}

	repos, err := handler.repositories.All()
	if err != nil {
		log.Printf("failed to get repos (%s) %s", err, req)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reposTemplate.Execute(w, repos)
	log.Printf("rendered repos/index with %d entries [%s]", len(repos), time.Since(startTime))
}