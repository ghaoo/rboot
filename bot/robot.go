package bot

type Robot struct {
	name string
	inMessage  chan Message
	outMessage chan Message
}

func New() *Robot {
	return &Robot{
		inMessage: make(chan Message),
		outMessage: make(chan Message),
	}
}

