package modules

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/filecopy"
	"github.com/gisvr/golib/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type TemplateReplace struct {
	Before []*Replace `yaml:"compile_before"`
	After  []*Replace `yaml:"compile_after"`
}

type Template struct {
	Base
	Variable `yaml:",inline"`
	Src      string           `yaml:"src"`
	Dst      string           `yaml:"dst"`
	Mode     string           `yaml:"mode"`
	Content  string           `yaml:"content"` //path为空是从这里读，content是一个key，具体内容在ctx的map里
	FVar     string           `yaml:"fvar"`
	Repl     *TemplateReplace `yaml:"replace"`
}

func GenUUID() string {
	// generate 32 bits timestamp
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}

func (this *Template) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	this.Base.Init(c, ctx)
	if this.Content == "" {
		//随机生成一个
		this.Content = GenUUID()
	}
	//加载预定义变量
	err := this.Variable.Init(this.c, this.ctx)
	if err != nil {
		return err
	}
	if this.FVar != "" {
		fvar, err := this.Compile(this.FVar)
		if err != nil {
			return err
		}
		err = compile.LoadVarsFromFile(fvar, this.c)
		if err != nil {
			return err
		}
	}
	//todo template的replace必须得一先init后才能运行，这里为了效率，作了优化，所有的replace一次pipeline执行
	if this.Repl != nil {
		for _, repl := range this.Repl.Before {
			repl.Content = this.Content
			err := repl.Init(this.c, this.ctx)
			if err != nil {
				return err
			}
		}
		for _, repl := range this.Repl.After {
			repl.Content = this.Content
			err := repl.Init(this.c, this.ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Template) Run() error {
	defer log.Infof("%s end\n\n", this.name)
	src := this.Src
	dst := this.Dst
	mode := this.Mode
	if src == "" {
		return errors.New("src is empty")
	}
	if dst == "" {
		return errors.New("dst is empty")
	}
	var err error
	src = strings.Replace(src, "\\", "/", -1)
	src, err = this.Compile(src)
	if err != nil {
		return err
	}
	dst = strings.Replace(dst, "\\", "/", -1)
	dst, err = this.Compile(dst)
	if err != nil {
		return err
	}
	if mode != "" {
		mode, err = this.Compile(mode)
		if err != nil {
			return err
		}
	} else {
		mode = "0755"
	}
	um, err := strconv.ParseInt(mode, 8, 0)
	if err != nil {
		return err
	}
	fmode := os.FileMode(um)
	log.Infof("%s.Run before copy. file mode: %s, parse: %d， fmode: %s", this.name, mode, um, fmode)
	return this.copy(src, dst, fmode)
}

func (this *Template) copy(src, dst string, mode os.FileMode) error {
	return filecopy.Copy(src, dst, mode, func(src, dst *os.File) (int64, error) {
		buf, err := ioutil.ReadAll(src)
		if err != nil {
			return 0, err
		}
		key, err := this.Compile(this.Content)
		if err != nil {
			return 0, err
		}
		content := new(string)
		*content = string(buf)
		this.SetKV(key, content)
		//compile before replace
		if this.Repl != nil && this.Repl.Before != nil {
			for _, repl := range this.Repl.Before {
				repl.SetKV(key, content)
				err = repl.Run()
				if err != nil {
					return 0, err
				}
			}
		}
		//compile
		*content, err = compile.BuildTemplateString(*content, this.c)
		if err != nil {
			return 0, err
		}
		//compile after replace
		if this.Repl != nil && this.Repl.After != nil {
			for _, repl := range this.Repl.After {
				repl.SetKV(key, content)
				err = repl.Run()
				if err != nil {
					return 0, err
				}
			}
		}
		len, err := dst.WriteString(*content)
		return int64(len), err
	})
}

func (this *Template) SetVar(key string, value interface{}) {
	this.Base.SetVar(key, value)
	if this.Repl != nil {
		for _, repl := range this.Repl.Before {
			repl.SetVar(key, value)
		}
		for _, repl := range this.Repl.After {
			repl.SetVar(key, value)
		}
	}
}
