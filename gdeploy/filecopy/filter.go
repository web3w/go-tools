package filecopy

import (
	"github.com/gisvr/golib/log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DirFilter(dir, filter string) ([]string, error) {
	var res []string
	if filter == "" {
		// 如果是空的filter，那么就搞所有文件
		filter = ".*"
	}
	re := regexp.MustCompile(filter)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		//windows dir fix
		path = strings.Replace(path, "\\", "/", -1)
		log.Infof("try filter file: %s", path)
		if re.MatchString(path) {
			res = append(res, path)
			log.Infof("get filter file: %s", path)
		}
		return nil
	})
	return res, err
}
