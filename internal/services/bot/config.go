package bot

type Config struct {
	Token string `env:"BOT_TOKEN"`
	Debug bool   `env:"BOT_DEBUG" default:"false"`
}
