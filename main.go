package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/concourse/go-concourse/concourse"
	"github.com/concourse/skymarshal/provider"
)

var (
	concourseTeam     = "main"
	concourseURL      = os.Getenv("CONCOURSE_URL")
	concourseUsername = os.Getenv("CONCOURSE_USERNAME")
	concoursePassword = os.Getenv("CONCOURSE_PASSWORD")
)

func main() {
	if concourseURL == "" {
		log.Fatal("environment variable required: ", concourseURL)
	}
	if concourseUsername == "" {
		log.Fatal("environment variable required: ", concourseUsername)
	}
	if concoursePassword == "" {
		log.Fatal("environment variable required: ", concoursePassword)
	}

	authToken := getAuthTokenForTeam(concourseURL, concourseTeam, concourseUsername, concoursePassword)
	concourseClient := concourse.NewClient(concourseURL, &http.Client{Transport: tokenTransport{authToken}}, false)

	workers, err := concourseClient.ListWorkers()
	if err != nil {
		log.Fatal("could not list workers: ", err)
	}

	for _, worker := range workers {
		if worker.State == "stalled" {
			fmt.Println("Pruning stalled worker", worker.Name)
			pruneErr := concourseClient.PruneWorker(worker.Name)
			if pruneErr != nil {
				log.Fatal("could not prune worker: ", pruneErr)
			}
		}
	}
}

func getAuthTokenForTeam(url, team, username, password string) provider.AuthToken {
	client := concourse.NewClient(url, &http.Client{Transport: basicAuthTransport{username, password}}, false)

	t := client.Team("main")
	authToken, err := t.AuthToken()
	if err != nil {
		log.Fatal("could not retreive auth token: ", err)
	}

	return authToken
}

type basicAuthTransport struct {
	username string
	password string
}

func (bat basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			bat.username, bat.password)))))
	return http.DefaultTransport.RoundTrip(req)
}

type tokenTransport struct {
	authToken provider.AuthToken
}

func (t tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", t.authToken.Type, t.authToken.Value))
	return http.DefaultTransport.RoundTrip(req)
}
