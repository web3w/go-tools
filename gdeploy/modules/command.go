package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Base
	Cmd      string   `yaml:"cmd"`
	Args     []string `yaml:"args"`
	Register string   `yaml:"register"`
}

func (this *Command) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	return this.Base.Init(c, ctx)
}

func (this *Command) Run() error {
	defer log.Infof("%s end\n\n", this.name)
	cmd, err := this.Compile(this.Cmd)
	if err != nil {
		return nil
	}
	if cmd == "" {
		return errors.New("cmd is empty")
	}
	args := make([]string, len(this.Args))
	for idx, arg := range this.Args {
		args[idx], err = this.Compile(arg)
		if err != nil {
			return err
		}
	}
	register, err := this.Compile(this.Register)
	if err != nil {
		return err
	}
	//执行命令
	log.Infof("%s %s", cmd, strings.Join(args, " "))
	if register != "" {
		err, stdoutmsg, stderrmsg := mCmdWithStdMsg(cmd, args...)
		if err != nil {
			log.Errorf("%s", stderrmsg)
			return err
		}
		if register != "" {
			this.SetKV(register, &stdoutmsg)
		}
	} else {
		r := exec.Command(cmd, args...)
		r.Stdout = os.Stdout
		r.Stderr = os.Stderr
		err = r.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func mCmdWithStdMsg(name string, args ...string) (error, string, string) {
	log.Infof("receive cmd: %s %v", name, args)

	/* #nosec */
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("bind to stderr err: %v", err), "", ""
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("bind to stderr err: %v", err), "", ""
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd start err: %v", err), "", ""
	}

	stdoutMsg, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Infof("read from stdout err: %v", err)
	}
	stderrMsg, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Infof("read from stderr err: %v", err)
	}

	err = cmd.Wait()
	return err, string(stdoutMsg), string(stderrMsg)
}
