package config

type AssetFilter struct {
	Match string `yaml:"match,omitempty"`
}

type AssetPuller struct {
	Name        string `yaml:"name"`
	Owner       string `yaml:"owner"`
	Repository  string `yaml:"repository"`
	Tag         string `yaml:"tag,omitempty"`
	Destination string `yaml:"destination"`
	Interval    uint64 `yaml:"interval,omitempty"`

	Filters []AssetFilter `yaml:"filters,omitempty"`
}

type Config struct {
	AssetPullers []AssetPuller `yaml:"asset-pullers"`
}
