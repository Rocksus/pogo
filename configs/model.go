package configs

type Config struct {
	Chat        ChatConfig
	Interpretor InterpretorConfig
	Port        string
	CertFile    string
	KeyFile     string
	Weather     WeatherConfig
}

type ChatConfig struct {
	ChannelAccessToken string
	ChannelSecret      string
	MasterID           string //ID of master user
}

type InterpretorConfig struct {
	Token string
}

type WeatherConfig struct {
	APIKey string
}
