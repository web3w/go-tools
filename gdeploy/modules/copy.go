package modules

import (
	"errors"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/filecopy"
	"github.com/gisvr/golib/log"
	"io"
	"os"
	"strconv"
	"strings"
)

type Copy struct {
	Base
	Src  string `yaml:"src"`
	Dst  string `yaml:"dst"`
	Mode string `yaml:"mode"`
}

func (this *Copy) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	this.Base.Init(c, ctx)
	return nil
}

func (this *Copy) Run() error {
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
	err = filecopy.Copy(src, dst, fmode, func(src, dst *os.File) (int64, error) {
		return io.Copy(dst, src)
	})
	return err
}
