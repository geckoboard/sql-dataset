package geckoboard

type Config struct {
	Key string
	URL string
}

func defaultConfig() Config {
	return Config{
		URL: "https://api.geckoboard.com",
	}
}

func (config *Config) mergeIn(other Config) {
	if other.Key != "" {
		config.Key = other.Key
	}

	if other.URL != "" {
		config.URL = other.URL
	}
}
