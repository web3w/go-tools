package compile

import (
	"github.com/flosch/pongo2"
	"github.com/gisvr/golib/config"
	"github.com/gisvr/golib/log"
	"gopkg.in/yaml.v2"
	"strings"
)

const (
	indentKey  = "indent1"
	encryptKey = "encrypt"
)

func indent(in string, count int) string {
	lines := strings.Split(in, "\n")
	if len(lines) <= 1 {
		return in
	}
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", count) + line
	}
	return strings.Join(lines, "\n")
}

func encrypt(in string) string {
	out, err := config.Encrypt(in)
	if err != nil {
		log.Fatalf("error encrypt: %v", err.Error())
	}
	return out
}

func indentFilterFn(in *pongo2.Value, params *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//对于map/array/slice得先格式化成字符串
	if !in.IsString() {
		if in.IsNil() {
			log.Fatalf("indent filter input is nil")
		}
		dat, err := yaml.Marshal(in.Interface())
		if err != nil {
			log.Fatalf("%v yam.marsha1 failed", in.Interface())
		}
		ss := string(dat)
		in = pongo2.AsValue(ss)
	}

	val := in.String()
	val = strings.Trim(val, " \t\r\n")
	cnt := params.Integer()
	out := indent(val, cnt)
	return pongo2.AsValue(out), nil
}

func encryptFilterFn(in *pongo2.Value, params *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()
	out := encrypt(val)
	return pongo2.AsValue(out), nil
}

func init() {
	pongo2.RegisterFilter(indentKey, indentFilterFn)
	pongo2.RegisterFilter(encryptKey, encryptFilterFn)
}
