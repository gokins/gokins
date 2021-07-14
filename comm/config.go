package comm

type Config struct {
	Server struct {
		Host      string   `yaml:"host"` //外网访问地址
		LoginKey  string   `yaml:"loginKey"`
		RunLimit  int      `yaml:"runLimit"`
		HbtpHost  string   `yaml:"hbtpHost"`
		Secret    string   `yaml:"secret"`
		Shells    []string `yaml:"shells"`
		DownToken string   `yaml:"DownToken"`
	} `yaml:"server"`
	Datasource struct {
		Driver string `yaml:"driver"`
		Url    string `yaml:"url"`
	} `yaml:"datasource"`
}
