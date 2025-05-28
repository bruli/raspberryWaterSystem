package telegram

const (
	HelpCommandName       CommandName = "help"
	StatusCommandName     CommandName = "status"
	LogCommandName        CommandName = "log"
	WaterCommandName      CommandName = "water"
	ActivateCommandName   CommandName = "activate"
	DeactivateCommandName CommandName = "deactivate"
	WeatherCommandName    CommandName = "weather"
	ZoneCommandName       CommandName = "zone"
)

type CommandName string

func (c CommandName) String() string {
	return string(c)
}

type Command struct {
	name                CommandName
	syntax, description string
}

func initCommands() [8]Command {
	return [8]Command{
		{HelpCommandName, "/help", "Show available commands"},
		{StatusCommandName, "/status", "Show current status"},
		{LogCommandName, "/log [limit]", "Show last log entries"},
		{WaterCommandName, "/water [zone] [seconds]", "Water zone for given seconds"},
		{ActivateCommandName, "/activate", "Activate server"},
		{DeactivateCommandName, "/deactivate", "Deactivate server"},
		{WeatherCommandName, "/weather", "Check current weather"},
		{ZoneCommandName, "/zone [id] [name] [relay1,relay2]", "Create a new zone"},
	}

}
