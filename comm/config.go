package comm

type Config struct {
	Server struct {
		Host     string `yaml:"host"` //外网访问地址
		Secret   string `yaml:"secret"`
		RunLimit int    `yaml:"run-limit"`
	} `yaml:"server"`
	Database struct {
		Driver string `yaml:"driver"`
		Url    string `yaml:"url"`
	} `yaml:"database"`
}
