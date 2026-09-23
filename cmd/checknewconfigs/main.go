package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/osbuild/blueprint/pkg/blueprint"
	"github.com/osbuild/image-builder/internal/common"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/distro/defs"
	"github.com/osbuild/image-builder/pkg/manifest"
	"github.com/osbuild/image-builder/pkg/manifestgen/manifestmock"
	"github.com/osbuild/image-builder/pkg/rpmmd"
)

type ConfigMapping struct {
	layers    []string
	distro    string
	arch      string
	imageType string
}

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

}

func main() {

	configMappings := []ConfigMapping{
		{
			layers: []string{
				"distro/fedora.yaml",
				"platform/x86_64-base.yaml",
				"type/qcow2.yaml",
			},
			distro:    "fedora-44",
			arch:      "x86_64",
			imageType: "generic-qcow2",
		},
		{
			layers: []string{
				"distro/fedora.yaml",
				"platform/x86_64-base.yaml",
				"type/ami.yaml",
			},
			distro:    "fedora-44",
			arch:      "x86_64",
			imageType: "ami",
		},
	}

	for idx, mapping := range configMappings {
		config, err := defs.Load(mapping.layers...)
		checkerr(err)

		fedora, err := defs.New(mapping.distro)
		checkerr(err)

		x86, err := fedora.GetArch(mapping.arch)
		checkerr(err)

		qcow2, err := x86.GetImageType(mapping.imageType)
		checkerr(err)

		converted, err := defs.ImageTypeFromConfig(config, qcow2)
		checkerr(err)

		basename := fmt.Sprintf("%s-%s-%s", mapping.distro, mapping.arch, mapping.imageType)

		jsonPrint(doManifest(converted), basename+".new.json")
		jsonPrint(doManifest(qcow2), basename+".old.json")
		fmt.Printf("[%02d] %s OK\n", idx, basename)
	}
}

func doManifest(imgType distro.ImageType) manifest.OSBuildManifest {
	options := distro.ImageOptions{}
	bp := new(blueprint.Blueprint)
	repos := []rpmmd.RepoConfig{
		{
			Name:     "example",
			BaseURLs: []string{"https://example.org/repo"},
		},
	}

	mf, _, err := imgType.Manifest(bp, options, repos, common.ToPtr(int64(99)))
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
