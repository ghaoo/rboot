package Rboot

// Response struct
type Response struct {
	Robot   *Robot
	Message *Message
	Match   []string
}

func NewResponse(robot *Robot, msg *Message) *Response {
	return &Response{
		Robot:   robot,
		Message: msg,
	}
}

func (res *Response) Text() string {

	return res.Message.Text
}

func (res *Response) Room() string {
	return res.Message.Room
}

func (res *Response) FromUser() string {
	return res.Message.FromUser
}

func (res *Response) ToUser() string {
	return res.Message.ToUser
}

func (res *Response) Send(strs ...string) error {
	for _, prov := range res.Robot.provider {
		go prov.Send(res, strs...)
	}
	return nil
}

func (res *Response) Reply(strs ...string) error {
	for _, prov := range res.Robot.provider {
		go prov.Reply(res, strs...)
	}
	return nil
}
