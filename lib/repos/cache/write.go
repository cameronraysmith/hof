package cache

import (
	"path/filepath"

	"github.com/go-git/go-billy/v5"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func Outdir(remote, owner, repo, tag string) string {
	outdir := filepath.Join(
		cacheBaseDir,
		remote,
		owner,
		repo+"@"+tag,
	)
	return outdir
}

func Write(remote, owner, repo, tag string, FS billy.Filesystem) error {
	outdir := Outdir(remote, owner, repo, tag)
	err := yagu.Mkdir(outdir)
	if err != nil {
		return err
	}
	return yagu.BillyWriteDirToOS(outdir, "/", FS)
}