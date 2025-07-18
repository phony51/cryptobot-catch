package config

type CatchConfig struct {
	Catcher    Credentials `json:"catcher"`
	Activator  Credentials `json:"activator"`
	Extractors Extractors  `json:"extractors"`
}

type Credentials struct {
	AppID    int    `json:"appID"`
	AppHash  string `json:"appHash"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
}
