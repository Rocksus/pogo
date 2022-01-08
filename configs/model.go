package configs

type Config struct {
	Chat        ChatConfig
	Interpretor InterpretorConfig
	Port        string
	CertFile    string
	KeyFile     string
	Weather     WeatherConfig
	News        NewsConfig
	Google      GoogleConfig
	MoneySheets MoneySheetsConfig
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

type NewsConfig struct {
	APIKey string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	Credentials  string
	Token string
}

type MoneySheetsConfig struct {
	SheetID string
}
