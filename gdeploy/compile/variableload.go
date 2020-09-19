package compile

import (
	"fmt"
	"github.com/gisvr/golib/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

var (
	indentRegexp   = regexp.MustCompile("([ ]+){{ *indent +([^ }]+) *}} *")
	encryptRegexp  = regexp.MustCompile("{{ *encrypt +([^\\.\"'][^\"' ]+) *}}")
	variableRegexp = regexp.MustCompile("{{ *([^\\.][^ ]+) *}}")
)

func extractStringMap(rf reflect.Value) map[string]interface{} {
	res := make(map[string]interface{})

	for _, key := range rf.MapKeys() {
		var k string
		switch f := key; key.Kind() {
		case reflect.Interface:
			value := f.Elem()
			if value.Kind() != reflect.String {
				return nil
			}
			k = value.String()
		case reflect.String:
			k = f.String()
		default:
			return nil
		}

		res[k] = rf.MapIndex(key).Interface()
	}

	return res
}

//根据环境提取对应的值
func extractVaule(k string, v interface{}, r *CompileT) (interface{}, error) {
	if sv, ok := v.(string); ok {
		return sv, nil
	}
	rf := reflect.ValueOf(v)
	var nv interface{}
	if r.Env != "" && rf.Kind() == reflect.Map {
		if sm := extractStringMap(rf); sm != nil {
			for _, mk := range []string{"." + r.Env, "default"} {
				if mv, ok := sm[mk]; ok {
					nv = mv
					break
				}
			}

			if nv != nil {
				for mk := range sm {
					if mk == "default" || mk[0] == '.' {
						continue
					}
					nv = nil
					log.Infof("will not use env for %v, because of key %v", k, mk)
					break
				}
			}
		}
	}

	if nv == nil {
		nv = v
	}
	return nv, nil
}

func yaml2str(i interface{}) string {
	if reflect.TypeOf(i).Kind() == reflect.String {
		str := reflect.ValueOf(i).String()
		return str
	}
	dat, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(dat)
	return strings.TrimRight(str, "\r\n")
}

func hasVariable(i interface{}) bool {
	s := yaml2str(i)
	return strings.Index(s, "{{") != -1
}

func internalFix(text string) string {
	res := text
	res = indentRegexp.ReplaceAllStringFunc(res, func(s string) string {
		fs := indentRegexp.FindStringSubmatch(s)
		log.Infof("capture indent: %v %#v", s, fs)
		return fmt.Sprintf("{{ %v|%s:%v }}", fs[2], indentKey, len(fs[1]))
	})

	res = encryptRegexp.ReplaceAllStringFunc(res, func(s string) string {
		fs := encryptRegexp.FindStringSubmatch(s)
		log.Infof("capture encrypt: %v %#v", s, fs)
		return fmt.Sprintf("{{ %v|%s }}", fs[1], encryptKey)
	})

	return res
}

func getValue(key string, r *CompileT) (interface{}, bool) {
	subkeys := strings.Split(key, ".")
	var subval interface{}
	subval = r.vars
	var ok bool
	for _, subkey := range subkeys {
		rf := reflect.ValueOf(subval)
		if rf.Kind() != reflect.Map {
			return nil, false
		}
		mval := extractStringMap(rf)
		subval, ok = mval[subkey]
		if !ok {
			return nil, false
		}
	}
	return subval, true
}

func variableReplace(str string, r *CompileT) (string, error) {
	var err error
	res := variableRegexp.ReplaceAllStringFunc(str, func(s string) string {
		fs := variableRegexp.FindStringSubmatch(s)
		key := fs[1]
		val, ok := getValue(key, r)
		if !ok {
			err = errors.New(fmt.Sprintf("%s", s))
			return s
		}

		val_as_string := reflect.ValueOf(val).String()
		//log.Infof("catture variable %s => %s", s, val_as_string)
		return val_as_string
	})
	if err != nil {
		return "", err
	}
	//log.Infof("%s => %s\n\n", str, res)
	return res, nil
}

func LoadVarsFromMap(localVars map[string]interface{}, r *CompileT) error {
	var ssMap = make(map[string]string)
	for k, v := range localVars {
		var sv interface{}
		sv, err := extractVaule(k, v, r)
		if err != nil {
			return err
		}

		if hasVariable(sv) {
			s := yaml2str(sv)
			s = internalFix(s)
			ssMap[k] = s
		} else {
			r.vars[k] = sv
		}
	}

	if len(ssMap) == 0 {
		return nil
	}
	i := 0
	maxLoop := 10
	log.Infof("replace variables begin")
	defer log.Infof("replace variable end")
	for ; i < maxLoop && len(ssMap) > 0; i++ {
		mapLen := len(ssMap)

		for k, tpl := range ssMap {
			//val, err := tpl.Execute(r.vars)
			val, err := variableReplace(tpl, r)
			if err != nil {
				//log.Warnf("%s not parse, try", err.Error())
				continue
			}
			log.Infof("replace{%s} %s => %s", k, tpl, val)
			r.vars[k] = val
			delete(ssMap, k)
		}

		if mapLen == len(ssMap) {
			return fmt.Errorf("can not resolve variables: %#v", ssMap)
		}
	}

	if i == maxLoop {
		return fmt.Errorf("can not resolve variables: %#v", ssMap)
	}
	//log.Infof("load vars success")
	return nil
}

func LoadVarsFromString(str string, r *CompileT) error {
	var localVars map[string]interface{}
	dat := []byte(str)
	err := yaml.Unmarshal(dat, &localVars)
	if err != nil {
		return errors.WithStack(err)
	}
	return LoadVarsFromMap(localVars, r)
}

func LoadVarsFromFile(fn string, r *CompileT) error {
	log.Infof("load vars from file{%s} begin", fn)
	defer log.Infof("load vars from file{%s} end", fn)
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	err = LoadVarsFromString(string(dat), r)
	if err != nil {
		log.Errorf("error load vars from file %v, err : %s", fn, err.Error())
		return err
	}
	return nil
}
