package notify

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"
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
	log.Debug().Msg("Loading last releases")
	for _, repository := range conf.Repositories {
		log.Debug().Msgf("Loading last release for %s", repository.Name)

		if repository.TagRegex != nil {
			if _, err := regexp.Compile(*repository.TagRegex); err != nil {
				log.Error().Err(err).Msg("Error compiling regex for " + repository.Name)
				panic(err)
			}
		}

		resp, err := http.Get("https://api.github.com/repos/" + repository.Name + "/releases/latest")
		if err != nil {
			log.Error().Err(err).Msg("Error getting latest release for " + repository.Name)
			panic(err)
		}
		defer resp.Body.Close()

		release := Release{}
		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing latest release for " + repository.Name)
			panic(err)
		}

		lastReleases[repository.Name] = release
	}

}

func Trigger() {
	var messages = map[string][]string{}

	for _, repository := range conf.Repositories {
		log.Debug().Msgf("Getting releases for %s", repository.Name)
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

			log.Info().Msgf("New release for %s: %s", repository.Name, release.TagName)
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
			log.Debug().Msgf("Sending %d messages to channel %s", len(msgs), channel)
			err := channelsMap[channel].SendMessages(msgs)
			if err != nil {
				log.Error().Err(err).Msg("Error while sending messages to " + channel)
			}
		}
	}
}

func getRepoReleases(name string) []Release {
	resp, err := http.Get("https://api.github.com/repos/" + name + "/releases")
	if err != nil {
		log.Error().Err(err).Msg("Error getting releases for " + name)
		return nil
	}
	defer resp.Body.Close()

	releases := []Release{}
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing releases for " + name)
		return nil
	}

	return releases
}
