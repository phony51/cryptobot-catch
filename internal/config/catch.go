package config

type CatchConfig struct {
	Catcher   Credentials `json:"catcher"`
	Activator Credentials `json:"activator"`
}

type Credentials struct {
	AppID    int    `json:"appID"`
	AppHash  string `json:"appHash"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
}
