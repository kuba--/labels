package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"
)

type (
	// Label ...
	Label struct {
		Name        string `json:"name,omitempty"`
		Color       string `json:"color,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// GithubClient ...
	GithubClient struct {
		*http.Client
	}
)

var (
	repo = flag.String("r", "", `Repository in format: <org>/<repo>`)
	file = flag.String("f", "default.json", "JSON file with labels")
	add  = flag.Bool("a", false, "Add a new label if doesn't exist")

	baseURL = func(repo string) string {
		return fmt.Sprintf("https://api.github.com/repos/%s/labels", repo)
	}
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <list | update> -r \"<org/repo>\" [options]\n", os.Args[0])
		fmt.Println("update options:")
		fmt.Printf("\t-f string\n\t\tJSON file with labels (default \"default.json\")\n")
		fmt.Printf("\t-a bool\n\t\tAdd a new label if doesn't exist (default false)\n")
		os.Exit(1)
	}
	if len(os.Args) < 3 {
		flag.Usage()
	}

	flag.CommandLine.Parse(os.Args[2:])
	client := NewGithubClient(os.Getenv("GITHUB_TOKEN"))
	switch cmd := os.Args[1]; cmd {
	case "list":
		if err := client.List(*repo); err != nil {
			log.Fatalln(err)
		}
	case "update":
		if err := client.Update(*repo, *file, *add); err != nil {
			log.Fatalln(err)
		}
	default:
		flag.Usage()
	}
}

// NewGithubClient ...
func NewGithubClient(token string) *GithubClient {
	return &GithubClient{oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)}
}

// List ...
func (gh *GithubClient) List(repo string) error {
	url := baseURL(repo)
	resp, err := gh.doRequest("GET", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var labels []*Label
	err = json.Unmarshal(body, &labels)
	if err != nil {
		return err
	}

	body, err = json.MarshalIndent(labels, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}

// Update ...
func (gh *GithubClient) Update(repo string, file string, add bool) error {
	var (
		data []byte
		err  error
	)
	if file == "" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(file)
	}
	if err != nil {
		return err
	}

	var labels []*Label
	err = json.Unmarshal(data, &labels)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, l := range labels {
		go func(lbl *Label) {
			wg.Add(1)
			defer wg.Done()

			var status string
			defer func() {
				fmt.Println("[" + lbl.Name + "]: " + status)
			}()

			url := fmt.Sprintf("%s/%s", baseURL(repo), lbl.Name)
			body, err := json.Marshal(lbl)
			if err != nil {
				status = err.Error()
			}

			resp, err := gh.doRequest("PATCH", url, bytes.NewBuffer(body))
			if err != nil {
				status = err.Error()
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound && add {
				url = baseURL(repo)
				resp, err = gh.doRequest("POST", url, bytes.NewBuffer(body))
				if err != nil {
					status = err.Error()
				}
				resp.Body.Close()
			}

			status = resp.Status

		}(l)
	}
	wg.Wait()

	return nil
}

func (gh *GithubClient) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", `application/vnd.github.symmetra-preview+json`)

	resp, err := gh.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
