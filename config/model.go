package config

// IConfig is an interface that helps you interact with the config module
type IConfig interface {
	Get() *Config
}

// IConfigModel used as the instance to the IConfig Interface
type IConfigModel struct {
	model *Config
}

//Config is a model that is used to pass the configuration through out the project
type Config struct {
	AppVersion string `yaml:"appVersion"`
	Server     Server `yaml:"server"`
	User       User   `yaml:"user"`
}

// Server contains server related configurations
type Server struct {
	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`
}

// HTTP contains http related configurations
type HTTP struct {
	Address string `yaml:"address"`
}

// GRPC contains GRPC related configurations
type GRPC struct {
	Address string `yaml:"address"`
}

// User contains user pkg specific config
type User struct {
	RatingsUrl    string `yaml:"ratingsUrl"`
	FavouritesUrl string `yaml:"favouritesUrl"`
}
