package command

type Command interface {
	CommandType() CommandType
}

type CommandType string
