package modules

import (
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCopy_Run(t *testing.T) {
	copy := &Copy{
		Src: "conf/var.yml",
		Dst: "testdir/aa/cc.yml",
	}
	c := compile.NewCompileT("dev")

	copy.Init(c, nil)
	err := copy.Run()
	assert.Equal(t, err, nil)

}
