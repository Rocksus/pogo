package configs

type Config struct {
	Chat        ChatConfig
	Interpretor InterpretorConfig
	Port        string
}

type ChatConfig struct {
	ChannelAccessToken string
	ChannelSecret      string
	MasterID           string //ID of master user
}

type InterpretorConfig struct {
	Token string
}
