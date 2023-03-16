package influx

type Options struct {
	Url    string `yaml:"url"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
	Token  string `yaml:"token"`
}

func Default() Options {
	return Options{}
}
