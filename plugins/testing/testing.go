package testing

import "rboot"

func setup(res *rboot.Response) error {
	switch res.Matcher {
	case `123`:
		res.Send(`1 or 2 or 3`)

	case `abc`:
		res.Send(`a or b or c`)

	}

	return nil
}

func init() {
	rboot.RegisterPlugin(`testing`, &rboot.Plugin{
		Action: setup,
		Ruleset: map[string]string{
			`123`: `1|2|3`,
			`abc`: `a|b|c`,
		},
	})
}
