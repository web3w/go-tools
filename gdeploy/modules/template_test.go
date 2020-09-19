package modules

import (
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTemplate_Run(t *testing.T) {
	//ckey:="ctx.t1"
	tobj := &Template{
		Src: "conf/t1.yml",
		Dst: "testdir/aa/t2.yml",
		//Content: ckey,
		FVar: "conf/var.yml",
		Repl: &TemplateReplace{
			Before: []*Replace{
				&Replace{
					Src: "hello4",
					Dst: "helloxxx",
				},
			},
			After: []*Replace{
				&Replace{
					Src: "helloxxx",
					Dst: "hello5",
				},
			},
		},
	}
	c := compile.NewCompileT("dev")
	ctx := NewContext()
	tobj.Init(c, ctx)
	err := tobj.Run()
	assert.Equal(t, err, nil)
}
