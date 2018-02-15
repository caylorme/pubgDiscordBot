package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/caylorme/pubgopgg"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var config = Config{
	Trigger: "!",
	Status:  "Rolyac",
}

type Config struct {
	Token   string
	Status  string
	Trigger string
}

func main() {
	fmt.Println("PUBG Discord Bot")
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		fmt.Println("Error decoding toml config", err.Error())
	} else {
		fmt.Println("Using Token: ", config.Token)
	}
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err.Error())
		return
	}
	dg.AddHandler(ready)
	dg.AddHandler(guildCreate)
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err.Error())
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("PUBG Discord Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, config.Status)

}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Bot Online.")
			return
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, config.Trigger) {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}
		params := strings.Split(m.Content, " ")
		switch command := strings.TrimPrefix(strings.Split(m.Content, " ")[0], config.Trigger); command {
		case "rank":
			if len(params) < 5 {
				return
			}
			pubgopgg, _ := pubgopgg.New()
			Player, err := pubgopgg.GetPlayer(params[1], params[2], params[3], params[4])
			if err != nil {
				_, _ = s.ChannelMessageSend(c.ID, err.Error())
			} else {
				_, _ = s.ChannelMessageSend(c.ID, Player.Username+" is ranked "+strconv.Itoa(Player.Ranks.Rating))
			}
		case "stats":
			if len(params) < 5 {
				return
			}
			pubgopgg, _ := pubgopgg.New()
			Player, err := pubgopgg.GetPlayer(params[1], params[2], params[3], params[4])
			if err != nil {
				_, _ = s.ChannelMessageSend(c.ID, err.Error())
			} else {
				_, _ = s.ChannelMessageSend(c.ID, "```Username: "+Player.Username+"  Rank: "+strconv.Itoa(Player.Ranks.Rating)+"\nRating: "+strconv.Itoa(Player.Stats.Rating)+" Matches: "+strconv.Itoa(Player.Stats.Matches_cnt)+" Wins: "+strconv.Itoa(Player.Stats.Win_matches_cnt)+" Top Tens: "+strconv.Itoa(Player.Stats.Topten_matches_cnt)+"\nKills: "+strconv.Itoa(Player.Stats.Kills_sum)+" Most Kills: "+strconv.Itoa(Player.Stats.Kills_max)+" Assists: "+strconv.Itoa(Player.Stats.Assists_sum)+" Headshots: "+strconv.Itoa(Player.Stats.Headshot_kills_sum)+" Deaths: "+strconv.Itoa(Player.Stats.Deaths_sum)+" Longest Kill: "+strconv.Itoa(Player.Stats.Longest_kill_max)+" Average Rank: "+strconv.FormatFloat(Player.Stats.Rank_avg, 'f', -1, 64)+" ADR: "+strconv.FormatFloat(Player.Stats.Damage_dealt_avg, 'f', -1, 64)+" Survival Average: "+strconv.FormatFloat(Player.Stats.Time_survived_avg, 'f', -1, 64)+"```")
			}

		default:
			_, _ = s.ChannelMessageSend(c.ID, command)
		}

		// Find the guild for that channel.
		// g, err := s.State.Guild(c.GuildID)
		// if err != nil {
		// 	// Could not find guild.
		// 	return
		// }

		// // Look for the message sender in that guild's current voice states.
		// for _, vs := range g.VoiceStates {
		// 	if vs.UserID == m.Author.ID {
		// 		err = playSound(s, g.ID, vs.ChannelID)
		// 		if err != nil {
		// 			fmt.Println("Error playing sound:", err)
		// 		}

		// 		return
		// 	}
		// }
	}
}
