package config

type BolterCfg struct {
	Logger   *Logger     `hcl:"logger,block"`
	Requests []*Requests `hcl:"requests,block"`
	Vegeta   *Vegeta     `hcl:"vegeta,block"`
}

type Logger struct {
	LogEnabled bool   `hcl:"log_enabled"`
	LoggerType int    `hcl:"logger_type"`
	FileName   string `hcl:"file_name"`
}
type Requests struct {
	Type    string   `hcl:"type,label"`
	Request *Request `hcl:"request,block"`
}

type Request struct {
	Jsonrpc    string   `hcl:"jsonrpc"`
	Method     string   `hcl:"method"`
	Parameters []string `hcl:"params"`
	Id         string   `hcl:"id"`
	HardCoded  bool     `hcl:"hard_coded"`
}
type Vegeta struct {
	Url      string  `hcl:"url"`
	Method   string  `hcl:"method"`
	IsPublic bool    `hcl:"is_public"`
	Rate     int     `hcl:"rate"`
	Duration int     `hcl:"duration"`
	Header   *Header `hcl:"header,block"`
}

// TODO: Add Basic (login/pass) auth

type Header struct {
	Auth   string `hcl:"auth"` //authorization type
	Bearer string `hcl:"bear"` //only Bearer token available at this time
}
