package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/mattn/go-zglob"
)

func (R *Runtime) AddDevTools(c *dagger.Container) (*dagger.Container) {
	c = c.WithExec([]string{"apt-get", "update", "-y"})
	c = c.WithExec([]string{"apt-get", "install", "-y", "tree"})
	return c
}

func (R *Runtime) ShowWorkdir(c *dagger.Container) *dagger.Container {
	c = c.WithExec([]string{"pwd"})
	c = c.WithExec([]string{"tree"})
	return c
}

func (R *Runtime) SetupTestingEnv(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {

	// add full code
	c = c.WithDirectory("/work", source)

	// set env vars
	c = c.WithEnvVariable("HOF_TELEMETRY_DISABLED", "1")
	c = c.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	c = c.WithEnvVariable("GITHUB_TOKEN_LEN", fmt.Sprintf("%d", len(os.Getenv("GITHUB_TOKEN"))))
	c = c.WithEnvVariable("HOF_FMT_VERSION", os.Getenv("HOF_FMT_VERSION"))

	// set fmt vars
	ver := "v0.6.10-rc.2" // todo, get this from the version used in the CLI command, consider falling back to some version if dirty or CI?
	c = c.WithEnvVariable("HOF_FMT_VERSION", ver)
	c = c.WithEnvVariable("HOF_FMT_HOST", "http://global-dockerd")

	// set testing dir
	c = c.WithWorkdir("/test")

	// uncomment at dev time if needed
	// c = R.AddDevTools(c)
	return c
}

func (R *Runtime) RunTestscriptDir(c *dagger.Container, source *dagger.Directory, name, dir, pattern string) error {

	d := source.Directory(dir)
	files, err := d.Entries(R.Ctx)
	if err != nil {
		return err
	}

	p := c

	// we want to run each as a separate fork of the testing container, in this way
	// each test gets a fresh environment and we can collect multiple errors before failing totally
	hadError := false
	for _, f := range files {
		ext := filepath.Ext(f)
		var F *dagger.File
		if pattern != "" {
			match, err := zglob.Match(pattern, f)		
			if err != nil {
				return err
			}
			if !match {
				continue
			}
			F = d.File(f)
		} else if ext == ".txt" || ext == ".txtar" {
			F = d.File(f)
		} else {
			continue
		}

		t := p

		t = t.WithMountedFile(filepath.Join("/test", f), F)
		t = t.WithExec([]string{"hof", "run", f})

		// now we only sync and check results once
		_, err = t.Sync(R.Ctx)
		if err != nil {
			hadError = true
		}
	}

	if hadError {
		return fmt.Errorf("errors while running %s in %s", name, dir)
	}

	return nil
}

func (R *Runtime) TestCommandFmt(c *dagger.Container, source *dagger.Directory) error {

	t := c

	t = t.WithExec([]string{"hof", "fmt", "start", "all"})
	t = t.WithExec([]string{"hof", "fmt", "info"})

	err := R.RunTestscriptDir(t, source, "test/fmt", "formatters/test", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestAdhocRender(c *dagger.Container, source *dagger.Directory) error {
	t := c

	err := R.RunTestscriptDir(t, source, "test/render", "test/render", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestMod(c *dagger.Container, source *dagger.Directory) error {
	t := c

	t = t.WithEnvVariable("GITLAB_TOKEN", os.Getenv("GITLAB_TOKEN"))
	t = t.WithEnvVariable("BITBUCKET_USERNAME", os.Getenv("BITBUCKET_USERNAME"))
	t = t.WithEnvVariable("BITBUCKET_PASSWORD", os.Getenv("BITBUCKET_PASSWORD"))

	err := R.RunTestscriptDir(t, source, "test/mod", "lib/mod/testdata", "")
	if err != nil {
		return err
	}

	t = t
	err = R.RunTestscriptDir(t, source, "test/mod/auth", "lib/mod/testdata/authd/apikeys", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestCreate(c *dagger.Container, source *dagger.Directory) error {
	p := c

	dirs := []string{
		"test/create/test_01",
		"test/create/test_02",
	}

	for _, dir := range dirs {
		d := source.Directory(dir)

		t := p
		t = t.WithDirectory("/test", d)
		t = t.WithExec([]string{"make", "test"})
		_, err := t.Sync(R.Ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (R *Runtime) TestStructural(c *dagger.Container, source *dagger.Directory) error {
	t := c

	err := R.RunTestscriptDir(t, source, "test/structural", "lib/structural", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestFlow(c *dagger.Container, source *dagger.Directory) error {
	p := c

	dirs := []string{
		"flow/testdata/bulk",
		"flow/testdata/tasks",
		"flow/testdata/tasks/api",
		"flow/testdata/tasks/db",
		"flow/testdata/tasks/ext",
		"flow/testdata/tasks/gen",
		"flow/testdata/tasks/hof",
		"flow/testdata/tasks/kv",
		"flow/testdata/tasks/os",
		"flow/testdata/tasks/st",
		"flow/testdata/concurrency",
	}

	for _, dir := range dirs {
		t := p
		err := R.RunTestscriptDir(t, source, dir, dir, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (R *Runtime) TestDatamodel(c *dagger.Container, source *dagger.Directory) error {
	t := c

	err := R.RunTestscriptDir(t, source, "test/datamodel", "lib/datamodel/test/testdata", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestCuecmd(c *dagger.Container, source *dagger.Directory) error {
	t := c

	err := R.RunTestscriptDir(t, source, "test/cuecmd", "lib/cuecmd/testdata", "")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestHack(c *dagger.Container, source *dagger.Directory) error {
	t := c

	err := R.RunTestscriptDir(t, source, "test/hack", "lib/cuecmd/testdata", "vet_*.txt")
	if err != nil {
		return err
	}

	return nil
}

