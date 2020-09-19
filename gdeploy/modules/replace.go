package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Replace struct {
	Base
	Path    string `yaml:"path"`    //读取写入同一个文件名
	Content string `yaml:"content"` //path为空是从这里读，content是一个key，具体内容在ctx的map里
	Src     string `yaml:"src"`     //被替换内容 oldtxt
	Regexp  string `yaml:"regexp"`  //正则替换 Src/Regexp只能有一个
	Dst     string `yaml:"dst"`     //替换内容  newtxt
}

func (this *Replace) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	this.Base.Init(c, ctx)
	this.Path = strings.Replace(this.Path, "\\", "/", -1)
	return nil
}

func (this *Replace) Run() error {
	defer log.Infof("%s end\n\n", this.name)
	if this.Path != "" && this.Content != "" {
		return errors.New("path and content only use one")
	}
	if this.Path == "" && this.Content == "" {
		return errors.New("path is empty")
	}
	path, err := this.Compile(this.Path)
	if err != nil {
		return err
	}
	contentkey, err := this.Compile(this.Content)
	if err != nil {
		return err
	}
	src, err := this.Compile(this.Src)
	if err != nil {
		return err
	}
	dst, err := this.Compile(this.Dst)
	if err != nil {
		return err
	}

	if path != "" {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		content := string(buf)
		content = this.replace(content, this.Regexp, src, dst)
		err = ioutil.WriteFile(path, []byte(content), stat.Mode())
		if err != nil {
			return err
		}
	} else if contentkey != "" {
		buf, err := this.GetKV(contentkey)
		if err != nil {
			return err
		}
		_, ok := buf.(*string)
		if !ok {
			return errors.New(fmt.Sprintf("%s not string", contentkey))
		}
		content := buf.(*string)
		*content = this.replace(*content, this.Regexp, src, dst)
	}
	return nil
}

func (this *Replace) replace(content, reg, src, dst string) string {
	if src == "" {
		if reg == "" {
			return content
		}
		replRegex := regexp.MustCompile(reg)
		return replRegex.ReplaceAllString(content, dst)
	}
	return strings.ReplaceAll(content, src, dst)
}
