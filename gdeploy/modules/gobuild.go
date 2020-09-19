package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"os"
	"os/exec"
)

type GoBuild struct {
	Base
	Pkgs   string   `yaml:"pkgs"`
	FBin   string   `yaml:"fbin"`
	Ver    bool     `yaml:"ver"`
	VerPkg string   `yaml:"verpkg"`
	VerFmt string   `yaml:"verfmt"` //版本的格式，由外部指定，可用变量gitversion/gitcommitid/gittreestate/builddate/platform
	Args   []string `yaml:"args"`   //编译参数
}

func (this *GoBuild) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	return this.Base.Init(c, ctx)
}

func (this *GoBuild) Run() error {
	defer log.Infof("%s end\n\n", this.name)
	pkgs := this.Pkgs
	var err error
	if pkgs != "" {
		pkgs, err = this.Compile(pkgs)
		if err != nil {
			return err
		}
	} else {
		pkgs = "." //默认为当前目录
	}
	fbin := this.FBin
	if fbin != "" {
		fbin, err = this.Compile(fbin)
		if err != nil {
			return err
		}
	} else {
		fbin = "." //默认为当前目录
	}
	args := make([]string, 0)
	args = append(args, "build")
	for _, arg := range this.Args {
		arg, err = this.Compile(arg)
		if err != nil {
			return err
		}
		args = append(args, arg)
	}
	if this.Ver {
		//id := getRepoCommitID()
		//ldflag := fmt.Sprintf("-X %s.GitSHA=%s -X %s.BuildTime=%s", this.VerPkg, id, this.VerPkg, time.Now().Format("2006-01-02T15:04:05"))
		if this.VerFmt == "" {
			return errors.New(fmt.Sprintf("verfmt is empty"))
		}
		ldflag, err := this.Compile(this.VerFmt)
		if err != nil {
			return err
		}
		args = append(args, "-ldflags", ldflag)
	}
	args = append(args, "-o", fbin, pkgs)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("%s", cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Infof("compile success!")
	return nil
}

func getRepoCommitID() string {
	out, err := exec.Command("git", "describe", "--dirty", "--always", "--tags").CombinedOutput()
	if err != nil {
		log.Warnf("git %v", err)
		out = []byte("------")
	}
	return string(out)

}
