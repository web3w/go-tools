package modules

import (
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/go-tools/gdeploy/shared"
	"github.com/prometheus/common/log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type GitRepo struct {
	TimeFmt  string `yaml:"timefmt"`
	Register bool   `yaml:"register"`
}

func (this *GitRepo) Init(c *compile.CompileT, ctx *Context) error {
	if !this.Register {
		return nil
	}

	mGitVer := map[string]string{}
	defer func(mVer map[string]string) {
		//for k, v := range mGitVer {
		//	c.SetVar(k, v)
		//}
		c.SetVar("gitrepo", mGitVer)
	}(mGitVer)
	timefmt := this.TimeFmt
	if timefmt == "" {
		timefmt = "2006-01-02T15:04:05"
	}
	buildDate := time.Now().Format(timefmt)
	mGitVer[shared.BuildDate] = buildDate

	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	mGitVer[shared.Platform] = platform

	ver, err := this.getBranch()
	if err != nil {
		return err
	}
	mGitVer[shared.GitVersion] = ver

	commitid, err := this.getCommitId()
	if err != nil {
		log.Errorf("%v", err)
		mGitVer[shared.GitCommitId] = "xxxxxxxx"
	} else {
		mGitVer[shared.GitCommitId] = commitid
	}

	commitdate, err := this.getCommitDate()
	if err != nil {
		log.Errorf("%v", err)
		mGitVer[shared.GitCommitDate] = "--------"
	} else {
		mGitVer[shared.GitCommitDate] = commitdate
	}
	return nil
}

func (this *GitRepo) GetBranch() (string, error) {
	return this.getBranch()
}

func (this *GitRepo) getBranch() (string, error) {
	cmd := exec.Command("git", []string{
		"branch",
		"-r",
		"--contains",
		"$(git rev-parse --abbrev-ref HEAD)",
	}...)
	log.Infof(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	rbranch := string(out)
	rbranch = strings.TrimPrefix(rbranch, "origin/")
	return rbranch, nil
}

func (this *GitRepo) getCommitId() (string, error) {
	cmd := exec.Command("git", []string{
		"rev-parse",
		"HEAD",
	}...)
	log.Infof(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	id := string(out)
	return id, nil
}

func (this *GitRepo) getCommitDate() (string, error) {
	cmd := exec.Command("git", []string{
		"log",
		"HEAD",
		"--pretty=format:\"%cd\"",
		"-1",
		fmt.Sprintf("--date=format:%Y%m%d"),
	}...)
	log.Infof(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	date := string(out)
	return date, nil
}
