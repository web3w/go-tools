package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/shared"
	"github.com/gisvr/golib/log"
	"gopkg.in/yaml.v2"
	"os/exec"
	"runtime"
	"time"
)

//支持go编译或script自定义编译
type BuildWrap struct {
	Base
	Variable `yaml:",inline"`
	TimeFmt  string       `yaml:"timefmt"`
	Prepare  *CommandWrap `yaml:"prepare"` //编译前准备
	GoBuild  *GoBuild     `yaml:"go"`
	Script   *CommandWrap `yaml:script`
	Post     *CommandWrap `yaml:"post"` //编译完成后清理
	runner   []IBase
}

func (this *BuildWrap) Init(c *compile.CompileT, ctx *Context) error {
	this.Base.Init(c, ctx)
	this.runner = make([]IBase, 0)
	//初始化之前先获取版本库版本信息
	err := this.setGitVersion()
	if err != nil {
		log.Errorf("%v", err)
	}
	//加载预定义变量
	err = this.Variable.Init(this.c, this.ctx)
	if err != nil {
		return err
	}
	//TODO prepare gobuild/script post 按顺序来
	if this.Prepare != nil {
		this.runner = append(this.runner, this.Prepare)
	}
	cc_num := 0
	if this.GoBuild != nil {
		this.runner = append(this.runner, this.GoBuild)
		cc_num++
	}
	if this.Script != nil {
		this.runner = append(this.runner, this.Script)
		cc_num++
	}
	if this.Post != nil {
		this.runner = append(this.runner, this.Post)
	}
	//只能一个编译命令
	if cc_num > 1 {
		out, _ := yaml.Marshal(this)
		return errors.New(fmt.Sprintf("[BuildWrap] synax error,\n%s", out))
	}
	for _, m := range this.runner {
		err := m.Init(this.c, this.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *BuildWrap) Run() error {
	for _, r := range this.runner {
		err := r.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *BuildWrap) SetVar(key string, value interface{}) {
	for _, r := range this.runner {
		r.SetVar(key, value)
	}
}

func (this *BuildWrap) setGitVersion() error {
	mGitVer := map[string]string{}
	defer func(mVer map[string]string) {
		for k, v := range mGitVer {
			this.c.SetVar(k, v)
		}
	}(mGitVer)
	timefmt := this.TimeFmt
	if timefmt == "" {
		timefmt = "2006-01-02T15:04:05"
	}
	buildDate := time.Now().Format(timefmt)
	mGitVer[shared.BuildDate] = buildDate

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	mGitVer[shared.Platform] = platform

	ver, err := this.Compile("{{ ctx.runenv.ref }}")
	if err != nil {
		return err
	}
	mGitVer[shared.GitVersion] = ver

	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		log.Errorf("%v", err)
		mGitVer[shared.GitCommitId] = "xxxxxxxx"
	} else {
		mGitVer[shared.GitCommitId] = string(out)
	}

	out, err = exec.Command("git", "log", "--pretty=format:\"%cd\"", "-1", "--date=format:%Y%m%d").CombinedOutput()
	if err != nil {
		log.Errorf("%v", err)
		mGitVer[shared.GitCommitDate] = "--------"
	} else {
		mGitVer[shared.GitCommitDate] = string(out)
	}
	return nil
}
