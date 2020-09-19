package modules

import (
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type IBase interface {
	Init(c *compile.CompileT, ctx *Context) error
	Run() error
	SetVar(key string, value interface{})
}

func getName(o interface{}) string {
	name := reflect.TypeOf(o).String()
	names := strings.Split(name, ".")
	name = names[len(names)-1]
	return strings.ToLower(name)
}

type Context struct {
	KV map[string]interface{}
}

func NewContext() *Context {
	return &Context{
		KV: make(map[string]interface{}),
	}
}

type Base struct {
	c    *compile.CompileT
	ctx  *Context
	name string //runner名称
}

func (this *Base) Init(c *compile.CompileT, ctx *Context) error {
	if this == nil {
		return errors.Errorf("[%s] not instince!!", this.name)
	}
	this.c = c.Clone()
	this.ctx = ctx
	return nil
}

func (this *Base) SetKV(key string, value interface{}) error {
	if this == nil {
		log.Fatalf("[%s] not instince!!", this.name)
	}
	if this.ctx == nil {
		log.Fatalf("[%s] not init", this.name)
	}
	this.ctx.KV[key] = value
	return nil
}

func (this *Base) GetKV(key string) (interface{}, error) {
	if this == nil {
		log.Fatalf("[%s] not instince!!", this.name)
	}
	if this.ctx == nil {
		log.Fatalf("[%s] not init", this.name)
	}
	if v, ok := this.ctx.KV[key]; !ok {
		return nil, errors.New(fmt.Sprintf("[%s] GetKV,key:%s not exist", this.name, key))
	} else {
		return v, nil
	}
}

func (this *Base) Compile(txt string) (string, error) {
	rlt, err := compile.BuildTemplateString(txt, this.c)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[%s] \"%s\" compile failed,err:%v", this.name, txt, err))
	}
	return rlt, nil
}

func (this *Base) SetVar(key string, value interface{}) {
	this.c.SetVar(key, value)
}
