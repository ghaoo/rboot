package bot

type Response struct {
	Robot   *Robot
	Message *Message
}

func NewResponse(bot *Robot, msg *Message) *Response {
	return &Response{
		Robot:   bot,
		Message: msg,
	}
}

func (res *Response) Send() {
	//
}

func (res *Response) Reply() {
	//
}
