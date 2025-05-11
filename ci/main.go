package main

import (
	"context"
	"dagger/ci/internal/dagger"
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

type JsonnetLib struct {
	Name string
	Dir  *dagger.Directory
}

func (m *Ci) WithJsonnetLibDirs(
	ctx context.Context,
	// +optional
	ctr *dagger.Container,
	libsRegex string,
) (jsonnetLibDirs []JsonnetLib, err error) {
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
		jsonnetLibDirs = append(jsonnetLibDirs, JsonnetLib{
			Name: filepath.Base(filepath.Dir(configFilepath)),
			Dir:  m.LibsRoot.Directory(filepath.Dir(configFilepath)),
		})
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
	jsonnetLibDirs, err := m.WithJsonnetLibDirs(ctx, nil, libsRegex)
	if err != nil {
		return nil, err
	}

	eg := errgroup.Group{}
	eg.SetLimit(parallel)

	artifactsDir := dag.Directory()
	for _, jsonnetLibDir := range jsonnetLibDirs {
		eg.Go(func() error {
			ctr := m.WithGoContainer()

			ctrInputDirLib := "/INPUT"
			ctrOutputDirRoot := "/OUTPUT"
			ctrOutputConfigDirRoot := "/OUTPUT_CONFIG"
			ctrOutputDirLib := filepath.Join(ctrOutputDirRoot, jsonnetLibDir.Name)

			ctr = ctr.
				WithMountedDirectory(ctrInputDirLib, jsonnetLibDir.Dir).
				WithExec([]string{
					"mkdir", "-p", ctrOutputDirLib,
				}).
				WithExec([]string{
					"jsonnet", "--create-output-dirs",
					"--multi", ctrOutputConfigDirRoot,
					"--jpath", ".",
					"--ext-str", ExtVarRepository + "=" + m.Repository,
					"--ext-str", ExtVarSiteUrl + "=" + m.SiteUrl,
					"--string", filepath.Join(ctrInputDirLib, m.ConfigFilenameJsonnet),
				}).
				WithExec([]string{
					"k8s-gen",
					"-o", ctrOutputDirLib,
					"-c", filepath.Join(ctrOutputConfigDirRoot, ConfigFilenameYAML),
				})

			outputDir := ctr.Directory(ctrOutputDirLib)
			artifactsDir = artifactsDir.WithDirectory(jsonnetLibDir.Name, outputDir)

			return nil
		})
	}

	return artifactsDir, eg.Wait()
}
