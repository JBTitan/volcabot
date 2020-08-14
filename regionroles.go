package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

func decodeRegionalIndicators(str string) string {
	var out string
	for _, r := range str {
		if r < 0x1F1E6 || r > 0x1F1FF {
			return ""
		}
		r -= 0x1F1E6 // 0=A
		r += 0x61    // 0x61=A
		out = out + string(r)
	}
	return out
}

func getOrCreateRegionRole(guildID string, name string) (*discordgo.Role, error) {
	role, err := getRegionRole(guildID, name)
	if err != nil {
		return role, err
	}

	if role == nil {
		role, err = discord.GuildRoleCreate(guildID)
		if err != nil {
			return nil, err
		}
		role, err = discord.GuildRoleEdit(guildID, role.ID, name, 0, false, 0, false)
		if err != nil {
			return nil, err
		}
	}
	return role, nil
}
func getRegionRole(guildID string, name string) (*discordgo.Role, error) {
	roles, err := discord.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	var role *discordgo.Role
	for i := range roles {
		if strings.HasPrefix(roles[i].Name, name) {
			role = roles[i]
			break
		}
	}
	return role, nil
}

func isRegionRoleMessage(channelID string, messageID string) (bool, error) {
	message, err := discord.ChannelMessage(channelID, messageID)
	if err != nil {
		return false, fmt.Errorf("error fetching message: %w", err)
	}
	me, err := discord.User("@me")
	if err != nil {
		return false, fmt.Errorf("error fetching @me: %w", err)
	}

	if message.Author.ID != me.ID {
		return false, nil
	}
	if message.Content != "React to this message with your home's flag!" {
		return false, nil
	}

	return true, nil
}

func init() {
	AddCommand("initregionreact", func(ctx *CommandContext) error {
		permissions, err := discord.UserChannelPermissions(ctx.User.ID, ctx.ChannelID)
		if err != nil {
			return err
		}
		if !(permissions&discordgo.PermissionAdministrator > 0) {
			return NewCommandError("You must have the Administrator permission to use that command")
		}
		_, err = discord.ChannelMessageSend(ctx.ChannelID, "React to this message with your home's flag!")
		return err
	})

	discord.AddHandler(func(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		log := log.WithFields(log.Fields{
			"message": reaction.MessageID,
			"channel": reaction.ChannelID,
			"guild":   reaction.GuildID,
			"user":    reaction.UserID,
			"emoji":   reaction.Emoji.Name,
		})

		isRegMessage, err := isRegionRoleMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			log.WithError(err).Errorln("error checking for region role message")
		}
		if !isRegMessage {
			return
		}

		reg := decodeRegionalIndicators(reaction.Emoji.Name)
		if len(reg) != 2 {
			return
		}
		log = log.WithField("regional_indicator", reg)

		role, err := getOrCreateRegionRole(reaction.GuildID, reaction.Emoji.Name)
		if err != nil {
			log.WithError(err).Errorln("error getting/creating region role")
			return
		}
		log = log.WithField("role", role.ID)

		err = discord.GuildMemberRoleAdd(reaction.GuildID, reaction.UserID, role.ID)
		if err != nil {
			log.WithError(err).Errorln("error applying role to user")
			return
		}

		log.Debug("added region role to user")
	})

	discord.AddHandler(func(s *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
		log := log.WithFields(log.Fields{
			"message": reaction.MessageID,
			"channel": reaction.ChannelID,
			"guild":   reaction.GuildID,
			"user":    reaction.UserID,
			"emoji":   reaction.Emoji.Name,
		})

		isRegMessage, err := isRegionRoleMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			log.WithError(err).Errorln("error checking for region role message")
		}
		if !isRegMessage {
			return
		}

		reg := decodeRegionalIndicators(reaction.Emoji.Name)
		if len(reg) != 2 {
			return
		}
		log = log.WithField("regional_indicator", reg)

		role, err := getRegionRole(reaction.GuildID, reaction.Emoji.Name)
		if err != nil {
			log.WithError(err).Errorln("error getting region role")
			return
		}
		if role == nil {
			log.Warnln("removed reaction for non-existent flag role")
			return
		}

		log = log.WithField("role", role.ID)

		err = discord.GuildMemberRoleRemove(reaction.GuildID, reaction.UserID, role.ID)
		if err != nil {
			log.WithError(err).Errorln("error removing region role from user")
		}
	})
}
