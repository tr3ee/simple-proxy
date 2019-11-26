package main

type Hook struct {
	hooks []func([]byte) []byte
}

func NewHook() *Hook {
	return &Hook{}
}

func (h *Hook) Run(p []byte) []byte {
	for _, hook := range h.hooks {
		p = hook(p)
	}
	return p
}

func (h *Hook) Add(hook func([]byte) []byte) {
	h.hooks = append(h.hooks, hook)
}
