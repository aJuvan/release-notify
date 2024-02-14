package notify

import (
	"github.com/aJuvan/release-notify/config"
	"github.com/aJuvan/release-notify/notify/channels"
	"github.com/aJuvan/release-notify/notify/channels/discord"
)

var conf *config.Config

var channelsMap = map[string]channels.Channel{}

func Setup(_conf *config.Config) {
	conf = _conf

	for _, discordNotif := range conf.Notifications.Discord {
		channelsMap[discordNotif.Name] = discord.NewChannel(discordNotif.Webhook)
	}

	initLastReleases()
}
