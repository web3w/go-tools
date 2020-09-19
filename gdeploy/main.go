package main

import (
	"flag"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/config"
	"github.com/gisvr/go-tools/gdeploy/modules"
	"github.com/gisvr/go-tools/gdeploy/shared"
	commsys "github.com/gisvr/golib/config"
	"github.com/gisvr/golib/log"
	"github.com/gisvr/golib/version"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	flag.StringVar(&commsys.Ref, "ref", "", "git branch")
	flag.Parse()
	version.Print()
	//加载参数
	cfg := config.Get()
	log.Init(cfg.Log)
	if commsys.ConfigFile == "" {
		arg0 := os.Args[0]
		arg0 = strings.Replace(arg0, "\\", "/", -1)
		arg0 = path.Dir(arg0)
		commsys.ConfigFile = arg0 + "/conf/gdeploy.yml"
	}
	if commsys.Ref == "" {
		commsys.Ref = os.Getenv("CI_COMMIT_REF_NAME")
		if commsys.Ref == "" {
			gitrepo := &modules.GitRepo{}
			commsys.Ref, _ = gitrepo.GetBranch()
		}
	}
	log.Infof("ref:%s, config: %s", commsys.Ref, commsys.ConfigFile)

	ref := commsys.Ref
	//创建运行时环境
	ctx := modules.NewContext()
	c := compile.NewCompileT(ref)
	//注册环境变量
	c.SetVar(shared.Ctx, ctx.KV)
	base := &modules.Base{}
	base.Init(c, ctx)
	base.SetKV(shared.SysEnv, cfg.SysEnvVar)
	fmtstr := cfg.TimeFmt
	if fmtstr == "" {
		fmtstr = "2006-01-02T15:04:05"
	}
	runEnvVar := &modules.RunEnvVar{
		Ref:         ref,
		CurrentTime: time.Now().Format(fmtstr),
	}
	runEnv := shared.Yaml_Struct2Interface(runEnvVar)
	base.SetKV(shared.RunEnv, runEnv)
	//加载工作文件
	tcfg, err := modules.LoadTargetConfig(cfg.ConfigFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if tcfg.Projects == nil {
		log.Infof("target is empty")
		return
	}
	for _, prj := range tcfg.Projects {
		//初始化
		log.Infof("project{%s} compile begin", prj.Name)
		err := prj.Init(c, ctx)
		if err != nil {
			log.Fatalf("project{%s} init failed, err: %+v", prj.Name, err)
			return
		}
		//执行
		err = prj.Run()
		if err != nil {
			log.Fatalf("project{%s} run failed, err: %+v", prj.Name, err)
			return
		}
		log.Infof("project{%s} compile finished!\n\n", prj.Name)
	}
	log.Infof("compile finished !")
}
