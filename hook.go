package rboot

type Hook interface {
	Fire(*Robot) error
}

type Hooks []Hook

func (hooks Hooks) Add(hook Hook) {
	hooks = append(hooks, hook)
}

func (hooks Hooks) Fire(bot *Robot) error {
	for _, hook := range hooks {
		if err := hook.Fire(bot); err != nil {
			return err
		}
	}

	return nil
}
