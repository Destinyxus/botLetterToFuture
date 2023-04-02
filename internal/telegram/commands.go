package telegram

type Commands struct {
	Start    bool
	Help     bool
	Stop     bool
	Goletter bool
}

func NewCommands() *Commands {
	return &Commands{
		Start:    false,
		Help:     false,
		Stop:     false,
		Goletter: false,
	}
}

func (c *Commands) CommandMode(cmd string) {
	switch cmd {
	case "start":
		c.Start = true
		c.Help = false
		c.Goletter = false
	case "help":
		c.Help = true
		c.Start = false
		c.Goletter = false
	case "goletter":
		c.Goletter = true
		c.Help = false
		c.Start = false
	case "reset":
		c.Goletter = false
	case "stop":
		c.Start = false
		c.Help = false
		c.Goletter = false
		c.Stop = false
	}
}
