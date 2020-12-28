package executor

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github-release-puller/config"

	"github.com/google/go-github/github"
)

func Start(cfg *config.Config) {
	client := github.NewClient(nil)

	for index, _ := range cfg.AssetPullers {
		go handlePuller(client, &cfg.AssetPullers[index])
	}
}

func handlePuller(client *github.Client, puller *config.AssetPuller) {
	log.Println("Initial puller", puller.Name)

	ticker := time.NewTimer(time.Microsecond)

	for range ticker.C {
		handlePull(client, puller)

		interval := puller.Interval
		if interval <= 0 {
			interval = 3600
		}

		ticker.Reset(time.Second * time.Duration(puller.Interval))
	}
}

func handlePull(client *github.Client, puller *config.AssetPuller) {
	log.Println("Pulling", puller.Name)

	var releases *github.RepositoryRelease
	var err error

	if puller.Tag == "" {
		releases, _, err = client.Repositories.GetLatestRelease(context.Background(), puller.Owner, puller.Repository)
	} else {
		releases, _, err = client.Repositories.GetReleaseByTag(context.Background(), puller.Owner, puller.Repository, puller.Tag)
	}

	if err != nil {
		log.Println("Pull", puller.Name, ":", err.Error())

		return
	}

	filters, err := compile(puller.Filters)
	if err != nil {
		log.Println("Compile regex for", puller.Name, ":", err.Error())

		return
	}

	for _, asset := range releases.Assets {
		if asset.Name == nil {
			continue
		}

		if filters.match(asset.GetName()) {
			output := path.Join(puller.Destination, asset.GetName())

			stat, err := os.Stat(output)
			if err != nil && !os.IsNotExist(err) {
				log.Println("Stat output file", output, ":", err.Error())

				continue
			}

			modified := time.Unix(0, 0)

			if err == nil {
				modified = stat.ModTime()
			}

			if asset.GetUpdatedAt().Before(modified) {
				log.Println(asset.GetName(), "up-to-date")

				continue
			}

			file, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
			if err != nil {
				log.Println("Write file", output, ":", err.Error())

				continue
			}

			response, err := http.Get(asset.GetBrowserDownloadURL())
			if err != nil {
				log.Println("Get", asset.GetBrowserDownloadURL(), ":", err.Error())

				file.Close()

				continue
			}

			if response.StatusCode/100 != 2 {
				body, err := ioutil.ReadAll(response.Body)
				if err == nil {
					err = errors.New(string(body))
				}

				log.Println("Get", asset.GetBrowserDownloadURL(), ":", err.Error())

				file.Close()
				response.Body.Close()

				continue
			}

			log.Println("Downloading", asset.GetBrowserDownloadURL())

			_, err = io.Copy(file, response.Body)
			if err == nil {
				log.Println(asset.GetName(), "pulled")
			}

			file.Close()
			response.Body.Close()

			if err != nil {
				log.Println("Pull", asset.GetName(), ":", err.Error())

				_ = os.Remove(output)
			}
		}
	}
}
