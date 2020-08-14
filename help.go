package main

import (
	"bytes"
	"fmt"
	"sort"
)

func init() {
	AddCommand("help", func(ctx *CommandContext) error {
		var text bytes.Buffer
		for i, section := range helpSections {
			fmt.Fprintf(&text, "__**%s**__\n", section.Name)
			for _, command := range section.Commands {
				fmt.Fprintf(&text, "    `%s` - %s\n", command.Command, command.Short)
			}
			if i < len(helpSections)-1 {
				fmt.Fprint(&text, "\n")
			}
		}
		return ctx.Reply(text.String())
	})
	AddHelpCommand("General", "help", "displays this message or detailed information about a command", "")
}

type helpSection struct {
	Name     string
	Commands []*helpCommand
}

func (section *helpSection) addCommand(cmd *helpCommand) {
	i := sort.Search(len(section.Commands), func(i int) bool {
		return section.Commands[i].Command >= cmd.Command
	})
	if i >= len(section.Commands) || section.Commands[i].Command != cmd.Command {
		section.Commands = append(section.Commands, nil)
		copy(section.Commands[i+1:], section.Commands[i:])
		section.Commands[i] = cmd
	} else {
		panic("command added to section twice")
	}
}

type helpCommand struct {
	Section string
	Command string
	Short   string
	Long    string
}

// helpSections is a sorted slice of existing help sections sorted by helpSection.Name
var helpSections []*helpSection

// helpCommands is a slice of existing help commands sorted by helpCommand.Command
var helpCommands []*helpCommand

func getOrCreateSection(sectionName string) *helpSection {
	i := sort.Search(len(helpSections), func(i int) bool {
		return helpSections[i].Name >= sectionName
	})
	if i >= len(helpSections) || helpSections[i].Name != sectionName {
		helpSections = append(helpSections, nil)
		copy(helpSections[i+1:], helpSections[i:])
		helpSections[i] = &helpSection{
			Name: sectionName,
		}
	}
	return helpSections[i]
}

func AddHelpCommand(sectionName string, cmd string, short string, long string) {
	command := &helpCommand{
		Section: sectionName,
		Command: cmd,
		Short:   short,
		Long:    long,
	}

	i := sort.Search(len(helpCommands), func(i int) bool {
		return helpCommands[i].Command >= cmd
	})
	if i >= len(helpCommands) || helpCommands[i].Command != cmd {
		helpCommands = append(helpCommands, nil)
		copy(helpCommands[i+1:], helpCommands[i:])
		helpCommands[i] = command
	} else {
		panic("command already added to helpCommands")
	}

	section := getOrCreateSection(sectionName)
	section.addCommand(command)
}
