package notify

import (
	"github.com/rs/zerolog/log"

	"github.com/aJuvan/release-notify/config"
	"github.com/aJuvan/release-notify/notify/channels"
	"github.com/aJuvan/release-notify/notify/channels/discord"
)

var conf *config.Config

var channelsMap = map[string]channels.Channel{}

func Setup(_conf *config.Config) {
	conf = _conf

	log.Debug().Msg("Loading notification channels")
	for _, discordNotif := range conf.Notifications.Discord {
		log.Debug().Msgf("Loading Discord channel %s", discordNotif.Name)
		channelsMap[discordNotif.Name] = discord.NewChannel(discordNotif.Webhook)
	}

	initLastReleases()
}
