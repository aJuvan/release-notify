package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Release struct {
	Name       string `json:"name"`
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Body       string `json:"body"`
}

func (r Release) String() string {
	return "# " + r.Name + "\n\n" + r.Body
}

var lastReleases = map[string]Release{}

func initLastReleases() {
	for _, repository := range conf.Repositories {
		if repository.TagRegex != nil {
			if _, err := regexp.Compile(*repository.TagRegex); err != nil {
				fmt.Println("Error compiling regex for " + repository.Name)
				panic(err)
			}
		}

		resp, err := http.Get("https://api.github.com/repos/" + repository.Name + "/releases/latest")
		if err != nil {
			fmt.Println("Error getting latest release for " + repository.Name)
			panic(err)
		}
		defer resp.Body.Close()

		release := Release{}
		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			fmt.Println("Error parsing latest release for " + repository.Name)
			panic(err)
		}

		lastReleases[repository.Name] = release
	}
}

func Trigger() {
	var messages = map[string][]string{}

	for _, repository := range conf.Repositories {
		var repoMessages = []string{}
		var releases = getRepoReleases(repository.Name)
		for _, release := range releases {
			if lastRelease, ok := lastReleases[repository.Name]; ok {
				if lastRelease.TagName == release.TagName {
					lastReleases[repository.Name] = releases[0]
					break
				}
			}

			if repository.TagRegex != nil {
				if match, _ := regexp.MatchString(*repository.TagRegex, release.TagName); !match {
					continue
				}
			}

			if !repository.Prerelease && release.Prerelease {
				break
			}

			repoMessages = append(repoMessages, release.String())
		}

		if len(repoMessages) > 0 {
			for i, j := 0, len(repoMessages)-1; i < j; i, j = i+1, j-1 {
				repoMessages[i], repoMessages[j] = repoMessages[j], repoMessages[i]
			}
		}

		for _, channel := range repository.NotificationChannels {
			if _, ok := messages[channel]; !ok {
				messages[channel] = []string{}
			}
			messages[channel] = append(messages[channel], repoMessages...)
		}
	}

	for channel, msgs := range messages {
		if len(msgs) > 0 {
			err := channelsMap[channel].SendMessages(msgs)
			if err != nil {
				fmt.Println("Error while sending messages to " + channel)
				fmt.Println(err)
			}
		}
	}
}

func getRepoReleases(name string) []Release {
	resp, err := http.Get("https://api.github.com/repos/" + name + "/releases")
	if err != nil {
		fmt.Println("Error getting releases for " + name)
		return nil
	}
	defer resp.Body.Close()

	releases := []Release{}
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		fmt.Println("Error parsing releases for " + name)
		return nil
	}

	return releases
}
