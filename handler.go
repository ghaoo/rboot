package Rboot

import (
	"regexp"
)

type handler interface {
	Handle(res *Response) error
}

func handlerMatch(r *regexp.Regexp, text string) bool {
	if !r.MatchString(text) {
		return false
	}
	return true
}

func handlerRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

// Handler type
type Handler struct {
	Pattern string
	Usage   string
	Run     func(res *Response) error
}

func (h *Handler) Handle(res *Response) error {
	switch {

	case h.Pattern == "":

		return h.Run(res)

	case h.match(res):

		res.Match = h.regexp().FindAllStringSubmatch(res.Message.Text, -1)[0]

		return h.Run(res)

	default:
		return nil
	}
}

func (h *Handler) regexp() *regexp.Regexp {
	return handlerRegexp(h.Pattern)
}

func (h *Handler) match(res *Response) bool {
	return handlerMatch(h.regexp(), res.Message.Text)
}
