package types

type Command struct {
	Command     string
	Aliases     []string
	Name        string
	Description string
	Usage       string
	Function    func([]string) error
}

func (c Command) MatchesCommand(command string) bool {
	if c.Command == command {
		return true
	}
	for _, alias := range c.Aliases {
		if alias == command {
			return true
		}
	}
	return false
}
