package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/osbuild/blueprint/pkg/blueprint"
	"github.com/osbuild/image-builder/internal/cmdutil"
	"github.com/osbuild/image-builder/internal/common"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/distro/defs"
	"github.com/osbuild/image-builder/pkg/manifest"
	"github.com/osbuild/image-builder/pkg/manifestgen/manifestmock"
	"github.com/osbuild/image-builder/pkg/rpmmd"
)

func checkerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}
}

func jsonPrint(thingie any, filename string) {
	c, err := json.MarshalIndent(thingie, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	_ = os.Remove(filename) // don't care if it fails

	if err := os.WriteFile(filename, c, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s OK\n", filename)
}

func main() {
	paths := []string{
		"distro/fedora.yaml",
		"platform/x86_64-base.yaml",
		"type/qcow2.yaml",
	}

	config, err := defs.Load(paths...)
	checkerr(err)

	// jsonPrint(config, "layered.json")

	fedora, err := defs.New("fedora-44")
	checkerr(err)

	x86, err := fedora.GetArch("x86_64")
	checkerr(err)

	qcow2, err := x86.GetImageType("generic-qcow2")
	checkerr(err)

	converted, err := defs.ImageTypeFromConfig(config, qcow2)
	checkerr(err)
	jsonPrint(doManifest(converted), "converted.json")
	jsonPrint(doManifest(qcow2), "old.json")
}

func doManifest(imgType distro.ImageType) manifest.OSBuildManifest {
	options := distro.ImageOptions{}
	seedArg, err := cmdutil.NewRNGSeed()
	bp := new(blueprint.Blueprint)
	repos := []rpmmd.RepoConfig{
		{
			Name:     "example",
			BaseURLs: []string{"https://example.org/repo"},
		},
	}

	mf, _, err := imgType.Manifest(bp, options, repos, &seedArg)
	checkerr(err)

	archName := imgType.Arch().Name()
	depsolvedSets, err := manifestmock.Depsolve(common.Must(mf.GetPackageSetChains()), archName, nil, false)
	checkerr(err)
	containerSpecs := manifestmock.ResolveContainers(mf.GetContainerSourceSpecs())
	commitSpecs := manifestmock.ResolveCommits(mf.GetOSTreeSourceSpecs())
	flatpakSpecs := manifestmock.ResolveFlatpaks(mf.GetFlatpakSourceSpecs())
	mfs, err := mf.Serialize(depsolvedSets, containerSpecs, commitSpecs, flatpakSpecs, nil)
	checkerr(err)
	return mfs
}
