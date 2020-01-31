package configs

type Config struct {
	Line LineConfig
	Wit	WitConfig
}

type LineConfig struct {
	ChannelAccessToken string
	ChannelSecret string
	MasterID string //ID of master user
}

type WitConfig struct {
	Token string
}