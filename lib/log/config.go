package log

type Config struct {
	Level         string `env:"LOG_LEVEL" envDefault:"debug"`
	IncludeCaller bool   `env:"LOG_INCLUDE_CALLER" envDefault:"false"`
	Format        string `env:"LOG_FORMAT" envDefault:"text"`
	Output        string `env:"LOG_OUTPUT" envDefault:"stdout"`
}
