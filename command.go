package main

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

// AddCommand adds a command with the given name to the bot. The function will be called whenever the command is executed.
// Upon execution, if an error fulfilling the interface CommandError is returned from CommandFunc, that error will be displayed to the user. If the error does not satisfy CommandError, a generic error message will be displayed instead. In all cases, the message that triggered the command will be deleted after CommandFunc is executed.
func AddCommand(cmd string, fn CommandFunc) {
	if commands == nil {
		commands = make(map[string]CommandFunc)
	}
	commands[cmd] = fn
}

type CommandFunc func(ctx *CommandContext) error

type CommandContext struct {
	Args []string

	ChannelID string
	Message   *discordgo.Message
	Member    *discordgo.Member // nil if sent in DMs
	User      *discordgo.User
}

func (ctx *CommandContext) Reply(content string) error {
	space := " "
	if strings.Contains(content, "\n") {
		space = "\n"
	}
	message := &discordgo.MessageSend{
		Content: ctx.User.Mention() + space + content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: []string{ctx.User.ID},
		},
	}

	_, err := discord.ChannelMessageSendComplex(ctx.ChannelID, message)
	return err
}

// ReplyComplex replies to the message used to trigger the command with the given discordgo.MessageSend. It prepends a Mention to the sender of the message to the Content string of the message and adds the sender to AllowedMentions.
func (ctx *CommandContext) ReplyComplex(msg *discordgo.MessageSend) error {
	return nil
	// TODO
}

func init() {
	discord.AddHandler(func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		if len(msg.WebhookID) > 0 || msg.Author.Bot {
			return
		}
		text := strings.TrimSpace(msg.Content)
		if !strings.HasPrefix(text, "volca!") {
			return
		}
		text = strings.TrimPrefix(text, "volca!")
		text = strings.TrimSpace(text)
		args := strings.Split(text, " ")
		if len(args) == 0 {
			return
		}
		cmd := args[0]
		args = args[1:]
		cmdfn := lookupCommand(cmd)
		if cmdfn == nil {
			return
		}

		ctx := &CommandContext{
			Args:      args,
			ChannelID: msg.ChannelID,
			Member:    msg.Member,
			User:      msg.Author,
		}

		log := log.WithFields(log.Fields{
			"command": cmd,
			"sender":  msg.Author.ID,
			"channel": msg.ChannelID,
			"guild":   msg.GuildID,
		})
		log.Debug("executing command")

		err := (*cmdfn)(ctx)
		if err != nil {
			switch err.(type) {
			case CommandError:
				log := log.WithField("error", err)
				log.Trace("command returned a CommandError; replying with message")
				err := ctx.Reply(err.Error())
				if err != nil {
					log.WithError(err).Errorln("error sending CommandError message")
					return
				}
			default:
				log.WithError(err).Errorln("error executing command")
				err := ctx.Reply("an error was encountered processing that command")
				if err != nil {
					log.WithError(err).Errorln("error sending error message message")
					return
				}
			}
		}
	})
}

// CommandError represents an error returned from a command that may be displayed to a user.
type CommandError interface {
	error
	_commanderror() // TODO(katie): add actual methods to CommandError (as needed) to differentiate it from regular errors
}

type genericCmdErr struct {
	Message string
}

func (err *genericCmdErr) _commanderror() {}
func (err *genericCmdErr) Error() string  { return err.Message }

func NewCommandError(msg string) CommandError {
	return &genericCmdErr{Message: msg}
}

var commands map[string]CommandFunc

func lookupCommand(cmd string) *CommandFunc {
	fn, ok := commands[cmd]
	if ok {
		return &fn
	}
	return nil
}
