package shared

import "gopkg.in/yaml.v2"

func Yaml_Struct2Interface(in interface{}) interface{} {
	buf, _ := yaml.Marshal(in)
	var out interface{}
	yaml.Unmarshal(buf, &out)
	return out
}

//
func Yaml_Interface2Struct(in, out interface{}) error {
	buf, _ := yaml.Marshal(in)
	return yaml.Unmarshal(buf, out)
}
