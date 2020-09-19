package compile

import (
	"github.com/flosch/pongo2"
)

// CompileT ...
type CompileT struct {
	Env  string
	vars pongo2.Context
}

// SetVar set kv into self
func (c *CompileT) SetVar(k string, v interface{}) {
	c.vars[k] = v
}

// SetVars set a lot of kv into self
func (c *CompileT) SetVars(vars map[string]interface{}) {
	for k, v := range vars {
		c.vars[k] = v
	}
}

// ResetVars reset all vars
func (c *CompileT) ResetVars(vars map[string]interface{}) {
	c.vars = make(map[string]interface{})
	c.SetVars(vars)
}

func (c *CompileT) Clone() *CompileT {
	nc := NewCompileT(c.Env)
	nc.ResetVars(c.vars)
	return nc
}

// NewCompileT returns a new pointer
func NewCompileT(env string) *CompileT {
	c := &CompileT{
		Env:  env,
		vars: pongo2.Context{},
	}

	return c
}

func BuildTemplateString(str string, r *CompileT) (string, error) {
	txt := internalFix(str)
	tpl, err := pongo2.FromString(txt)
	if err != nil {
		return "", err
	}
	rlt, err := tpl.Execute(r.vars)
	//log.Infof(`%s => %s`, txt,rlt)
	return rlt, err
}
