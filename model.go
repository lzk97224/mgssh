package main

type HostConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Key  string `json:"key"`
}

func (h *HostConfig) Dail() error {
	if len(h.Pass) >= 1 {
		return dialSShWithPassword(h.Host, h.Port, h.User, h.Pass)
	}
	return dialSSHUseCommand(h.Host, h.Port, h.User, h.Key)
}
