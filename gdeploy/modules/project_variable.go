package modules

type ModuleVar struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
	Ref  string `yaml:"ref"`
}

type ProjectVar struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
	Ref  string `yaml:"ref"`
}

type SystemVar struct {
	SupportEnv []string `yaml:"supportenv"`
}

//运行时环境变量
type RunEnvVar struct {
	Ref         string `yaml:"ref"`
	Label       string `yaml:"label"`       //
	CurrentTime string `yaml:"currenttime"` //程序运行时的时间
}
