package dagger

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// so we don't have to pass these around everywhere
type Runtime struct {
	Ctx    context.Context
	Client *dagger.Client
}

var goVer = "golang:1.23"

func (R *Runtime) GolangImage(platform string) (*dagger.Container) {
	c := R.Client.
		Container(dagger.ContainerOpts{Platform: dagger.Platform(platform)}).
		From(goVer)

	// setup mod cache
	modCache := R.Client.CacheVolume(fmt.Sprintf("gomod-%s-%s", goVer, platform))
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// setup build cache
	buildCache := R.Client.CacheVolume(fmt.Sprintf("go-build-%s-%s", goVer, platform))
	c = c.WithMountedCache("/root/.cache/go-build", buildCache)

	// setup workdir
	c = c.WithWorkdir("/work")

	return c
}

func (R *Runtime) RuntimeContainer(builder *dagger.Container, platform string) (*dagger.Container) {
	hof := builder.File("hof")

	c := R.GolangImage(platform)
	c = c.WithFile("/usr/local/bin/hof", hof)
	
	return c
}

func (R *Runtime) FetchDeps(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {

	// get deps
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "download"})

	// c = c.WithDirectory("/work", source)
	return c
}

func (R *Runtime) BuildHof(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {

	// exclude files we don't need so we can avoid cache misses?
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"changelogs",
			"ci",
			"docs",
			"hack",
			"images",
			"notes",
			"test", 
		},
	})

	c = c.WithEnvVariable("CGO_ENABLED", "0")

	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	return c
}

func (R *Runtime) BuildHofMatrix(c *dagger.Container) (*dagger.Directory, error) {
	// the matrix
	geese := []string{"linux", "darwin"}
	goarches := []string{"amd64", "arm64"}

	outputs := R.Client.Directory()

	// build matrix for writing to host
	for _, goos := range geese {
		for _, goarch := range goarches {
			// create a directory for each OS and architecture
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// set local env vars
			build := c.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch)

			// run the build
			build = build.WithExec([]string{"go", "build", "-o", path, "./cmd/hof"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	return outputs, nil
}

func (R *Runtime) HofVersion(c *dagger.Container) error {
	t := c.WithExec([]string{"hof", "version"})

	_, err := t.Sync(R.Ctx)
	return err
}
