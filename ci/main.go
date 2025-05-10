// A generated module for Ci functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ci/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	DefaultCtrImgFind  = "docker.io/alpine:3.20"
	ConfigFilenameYAML = "config.yml" // see: jsonnet/config.jsonnet

	ExtVarRepository = "repository"
	ExtVarSiteUrl    = "siteUrl"
)

type Ci struct {
	RepoRoot              *dagger.Directory
	LibsRoot              *dagger.Directory
	ConfigFilenameJsonnet string
	Repository            string
	SiteUrl               string
	mu                    sync.Mutex
}

func New(
	// +defaultPath="/"
	repoRoot *dagger.Directory,
	// +defaultPath="/"
	libsRoot *dagger.Directory,
	// +default="config.jsonnet"
	configFilenameJsonnet string,
	// +default="github.com/jsonnet-libs/"
	repository string,
	// +default="https://jsonnet-libs.github.io/"
	siteUrl string,
) *Ci {
	repoRoot = repoRoot.Filter(dagger.DirectoryFilterOpts{
		Exclude: []string{
			"ci/**/*",
			"gen/**/*",
		},
	})
	if !strings.HasSuffix(repository, "/") {
		repository += "/"
	}
	if !strings.HasSuffix(siteUrl, "/") {
		siteUrl += "/"
	}

	return &Ci{
		RepoRoot:              repoRoot,
		LibsRoot:              libsRoot,
		ConfigFilenameJsonnet: configFilenameJsonnet,
		Repository:            repository,
		SiteUrl:               siteUrl,
		mu:                    sync.Mutex{},
	}
}

func (m *Ci) WithGoContainer() *dagger.Container {
	ctr := dag.Container().
		Build(m.RepoRoot)
	return ctr
}

func (m *Ci) WithJsonnetLibDirs(
	ctx context.Context,
	// +optional
	ctr *dagger.Container,
	libsRegex string,
) (jsonnetLibDirs []*dagger.Directory, err error) {
	if ctr == nil {
		ctr = dag.Container().From(DefaultCtrImgFind)
	}

	ctr = ctr.
		WithMountedDirectory("/WORK", m.LibsRoot).
		WithWorkdir("/WORK")

	findCmd := []string{
		"find",
		".",
		"-type", "f",
		"-name", m.ConfigFilenameJsonnet,
	}
	var shellCmd []string
	if libsRegex != "" {
		grepCmd := []string{
			"grep", "-E", libsRegex,
		}

		shellCmd = []string{"sh", "-c",
			strings.Join(findCmd, " ") + " | " + strings.Join(grepCmd, " "),
		}
	} else {
		shellCmd = []string{"sh", "-c",
			strings.Join(findCmd, " "),
		}
	}

	configFileOutput, err := ctr.
		WithExec(shellCmd).
		Stdout(ctx)
	if err != nil {
		return
	}

	for _, configFilepath := range strings.Split(strings.TrimSpace(configFileOutput), "\n") {
		jsonnetLibDirs = append(jsonnetLibDirs, m.LibsRoot.Directory(filepath.Dir(configFilepath)))
	}

	return
}

func (m *Ci) Build(
	ctx context.Context,
	// +optional
	libsRegex string,
	// +default=5
	parallel int,
) (*dagger.Directory, error) {
	ctr := m.WithGoContainer()

	jsonnetLibDirs, err := m.WithJsonnetLibDirs(ctx, nil, libsRegex)
	if err != nil {
		return nil, err
	}

	eg, gctx := errgroup.WithContext(ctx)
	eg.SetLimit(parallel)

	artifactsDir := dag.Directory()
	for _, jsonnetLibDir := range jsonnetLibDirs {
		eg.Go(func() error {
			ctrInputDirLib := "/INPUT"
			ctrOutputDirRoot := "/OUTPUT"

			ctrWorkdirOriginal, err := ctr.Workdir(ctx)
			if err != nil {
				return err
			}

			ctr = ctr.WithMountedDirectory(ctrInputDirLib, jsonnetLibDir)

			jsonnetLibName, err := ctr.
				WithWorkdir(ctrInputDirLib).
				WithExec([]string{
					"sh", "-c",
					fmt.Sprintf(`cat %s | grep 'name=' | cut -d"'" -f2`, m.ConfigFilenameJsonnet),
				}).
				Stdout(gctx)
			if err != nil {
				return err
			}
			jsonnetLibName = strings.TrimSpace(jsonnetLibName)
			if jsonnetLibName == "" {
				return fmt.Errorf("could not find name for config file: '%s'", m.ConfigFilenameJsonnet)
			}

			ctrOutputDirLib := filepath.Join(ctrOutputDirRoot, jsonnetLibName)

			ctr = ctr.
				WithWorkdir(ctrWorkdirOriginal).
				WithExec([]string{
					"mkdir", "-p", ctrOutputDirLib,
				}).
				WithExec([]string{
					"jsonnet", "--create-output-dirs",
					"--multi", ctrOutputDirLib,
					"--jpath", ".",
					"--ext-str", ExtVarRepository + "=" + m.Repository,
					"--ext-str", ExtVarSiteUrl + "=" + m.SiteUrl,
					"--string", filepath.Join(ctrInputDirLib, m.ConfigFilenameJsonnet),
				}).
				WithExec([]string{
					"k8s-gen",
					"-o", ctrOutputDirLib,
					"-c", filepath.Join(ctrOutputDirLib, ConfigFilenameYAML),
				})

			outputDir := ctr.Directory(ctrOutputDirLib)
			artifactsDir = artifactsDir.WithDirectory(jsonnetLibName, outputDir)

			return nil
		})
	}

	return artifactsDir, eg.Wait()
}
