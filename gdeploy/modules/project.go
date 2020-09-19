package modules

import (
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/shared"
	"github.com/gisvr/golib/log"
	"github.com/pkg/errors"
	"strings"
)

type Project struct {
	Base
	Variable   `yaml:",inline"`
	Labels     []string            `yaml:"labels"` //暂时没什么用
	Name       string              `yaml:"name"`   //模块名称
	Conditions []*Condition        `yaml:"condtions"`
	Type       string              `yaml:"type"` //模块类型 比如bin/java/fileraw
	FVar       string              `yaml:"fvar"`
	Build      *BuildWrap          `yaml:"build"` //编译 目前只支持go
	Files      []*InstructionsWrap `yaml:"files"` //文件操作
	Shell      []*CommandWrap      `yaml:"shell"`
	Modules    []*Module           `yaml:"modules"`
	Order      []*OrderWrap        `yaml:"order"` //有序命令序列
}

func (this *Project) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s{%s} begin", this.name, this.Name)
	err := this.Base.Init(c, ctx)
	if err != nil {
		return err
	}
	runEnv, err := this.GetKV(shared.RunEnv)
	if err != nil {
		return err
	}
	var runEnvVar RunEnvVar
	err = shared.Yaml_Interface2Struct(runEnv, &runEnvVar)
	if err != nil {
		return errors.Errorf("set runenv to ctx")
	}

	var cc *Condition
	if this.Conditions != nil {
		for _, cond := range this.Conditions {
			if cond.Ref == runEnvVar.Ref {
				cc = &Condition{
					Ref: cond.Ref,
					Env: cond.Env,
				}
				break
			}
		}
		if cc == nil {
			//当前分支不部署
			log.Infof("%s{%s} ref:%s don't deploy", this.name, this.Name, runEnvVar.Ref)
			return nil
		}
	}
	//如果没有设cond，则env=ref
	if cc == nil {
		cc = &Condition{
			Ref: runEnvVar.Ref,
			Env: runEnvVar.Ref,
		}
	}

	sysEnv, err := this.GetKV(shared.SysEnv)
	if err != nil {
		return err
	}
	var sysEnvVar SystemVar
	err = shared.Yaml_Interface2Struct(sysEnv, &sysEnvVar)
	if err != nil {
		return err
	}

	//是否是支持有env
	if sysEnvVar.SupportEnv != nil {
		deploy := false
		for _, env := range sysEnvVar.SupportEnv {
			if env == cc.Env {
				deploy = true
				break
			}
		}
		if !deploy {
			log.Infof("%s{%s} env:%s not support", this.name, this.Name, cc.Env)
			return nil
		}
	}

	prjEnvVar := &ProjectVar{
		Name: this.Name,
		Ref:  cc.Ref,
		Env:  cc.Env,
	}
	prjEnv := shared.Yaml_Struct2Interface(prjEnvVar)
	this.SetKV(shared.PrjEnv, prjEnv)
	//环境变量已经出来了，加载参数之前需要更改
	this.c.Env = cc.Env

	//加载预定义变量
	err = this.Variable.Init(this.c, this.ctx)
	if err != nil {
		return err
	}
	if this.FVar != "" {
		fvar := strings.Replace(this.FVar, "\\", "/", -1)
		fvar, err := this.Compile(fvar)
		if err != nil {
			return err
		}
		//加载vars
		err = compile.LoadVarsFromFile(fvar, this.c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Project) Run() error {
	defer log.Infof("%s{%s} end\n\n", this.name, this.Name)
	if this.Build != nil {
		err := this.Build.Init(this.c, this.ctx)
		if err != nil {
			return err
		}
		err = this.Build.Run()
		if err != nil {
			return err
		}
	}
	for _, file := range this.Files {
		err := file.Init(this.c, this.ctx)
		if err != nil {
			return err
		}
		err = file.Run()
		if err != nil {
			return err
		}
	}
	for _, shell := range this.Shell {
		err := shell.Init(this.c, this.ctx)
		if err != nil {
			return err
		}
		err = shell.Run()
		if err != nil {
			return err
		}
	}

	//project 执行完成后才开始module的初始化执行
	for _, module := range this.Modules {
		if module.Conditions == nil {
			//module没有设置条件，就和project一样的条件
			module.Conditions = this.Conditions
		}
		err := module.Init(this.c, this.ctx)
		if err != nil {
			return err
		}
		err = module.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Project) SetVar(key string, value interface{}) {
	this.Base.SetVar(key, value)
	if this.Build != nil {
		this.Build.SetVar(key, value)
	}
	for _, file := range this.Files {
		file.SetVar(key, value)
	}
	for _, shell := range this.Shell {
		shell.SetVar(key, value)
	}

	for _, module := range this.Modules {
		module.SetVar(key, value)
	}
}
