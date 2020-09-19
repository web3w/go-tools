package filecopy

import (
	"errors"
	"fmt"
	"github.com/gisvr/golib/log"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type CopyFunc func(src, dst *os.File) (int64, error)
type FileCopy struct {
	cp CopyFunc
}

func (this *FileCopy) copydir(src, dst string, mode os.FileMode) error {
	fs, err := DirFilter(src, "")
	if err != nil {
		return err
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}

	src = strings.TrimRight(src, "/\\")
	for _, fn := range fs {
		srcfile, err := filepath.Abs(fn)
		if err != nil {
			return err
		}
		filename := srcfile[len(src):]
		filename = strings.Replace(filename, "\\", "/", -1)
		filename = strings.TrimLeft(filename, "/\\")
		//filename := fn[len(src):]
		err = this.copyfile(fn, path.Join(dst, filename), mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *FileCopy) setprop(src, dst string) error {
	st, err := os.Stat(src)
	if err != nil {
		return nil
	}
	return os.Chtimes(dst, time.Now(), st.ModTime())
}

func (this *FileCopy) copyfile(src, dst string, mode os.FileMode) error {
	if err := this.copyfilecontext(src, dst, mode); err != nil {
		return nil
	}
	return this.setprop(src, dst)
}

func (this *FileCopy) copyfilecontext(src, dst string, mode os.FileMode) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return errors.New(fmt.Sprintf("%s open failed,err=%v", src, err))
	}
	defer srcfile.Close()
	dir, _ := path.Split(dst)
	err = os.MkdirAll(dir, mode)
	if err != nil {
		return err
	}
	dstfile, err := os.Create(dst)
	if err != nil {
		return errors.New(fmt.Sprintf("%s create failed,err=%v", dst, err))
	}
	defer dstfile.Close()

	// copy origin mode, to avoid script file can't be execute.
	stat, err := os.Stat(src)
	if err != nil {
		return errors.New(fmt.Sprintf("%s stat read failed,err=%v", src, err))
	}
	err = dstfile.Chmod(stat.Mode())
	if err != nil {
		return errors.New(fmt.Sprintf("%s set file mode(%s) failed, err=%v", dst, stat.Mode().String(), err))
	}

	len, err := this.cp(srcfile, dstfile)
	if err != nil {
		return err
	}
	log.Infof("copy %s to %s,length=%d", src, dst, len)
	return nil
}

func (this *FileCopy) copy(src, dst string, mode os.FileMode) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		if src[len(src)-1] != '/' {
			//如果源目录没有/ 那目标是 源目录的最后一段复制过去:
			// cp a b => b/a
			// cp a b/ => b/a
			// cp a/ b => b/
			dst = path.Join(dst, path.Base(src))
		}
		return this.copydir(src, dst, mode)
	} else {
		dststat, err := os.Stat(dst)
		if err != nil {
			//dst不存在 最后一个是/就当目录，否则当文件名
			if dst[len(dst)-1] == '/' {
				// 如果目标是目录，那目标是原来的目标+源的最后一段 : cp a b/ => b/a
				dst = path.Join(dst, path.Base(src))
			}
		} else {
			//dst是目录
			if dststat.IsDir() {
				dst = path.Join(dst, path.Base(src))
			}
		}

		return this.copyfile(src, dst, mode)
	}
}

func Copy(src, dst string, mode os.FileMode, ccf CopyFunc) error {
	obj := &FileCopy{
		cp: ccf,
	}
	if obj.cp == nil {
		obj.cp = func(srcfile, dstfile *os.File) (int64, error) {
			return io.Copy(dstfile, srcfile)
		}
	}
	return obj.copy(src, dst, mode)
}
