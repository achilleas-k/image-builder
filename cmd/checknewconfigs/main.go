package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/osbuild/blueprint/pkg/blueprint"
	"github.com/osbuild/image-builder/internal/common"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/distro/defs"
	"github.com/osbuild/image-builder/pkg/manifest"
	"github.com/osbuild/image-builder/pkg/manifestgen/manifestmock"
	"github.com/osbuild/image-builder/pkg/rpmmd"
)

const diffpath = "./diff/"

type ConfigMapping struct {
	layers    []string
	distro    string
	arch      string
	imageType string
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}
}

func jsonMarshal(thingie any) []byte {
	c, err := json.MarshalIndent(thingie, "", "  ")
	check(err)
	return c
}

func save(data []byte, filename string) {
	os.Remove(filename)
	check(os.WriteFile(filename, data, 0600))
}

func compare(converted, orig manifest.OSBuildManifest, name string) {
	fmt.Printf("Checking %s: ", name)
	if slices.Compare(converted, orig) == 0 {
		fmt.Println("OK")
		return
	}

	fmt.Printf("DIFFER -> Saving to %s\n", diffpath)

	os.MkdirAll(diffpath, 0700)
	save(jsonMarshal(converted), filepath.Join(diffpath, name+".new.json"))
	save(jsonMarshal(orig), filepath.Join(diffpath, name+".old.json"))
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

	for _, mapping := range configMappings {
		config, err := defs.Load(mapping.layers...)
		check(err)

		fedora, err := defs.New(mapping.distro)
		check(err)

		x86, err := fedora.GetArch(mapping.arch)
		check(err)

		qcow2, err := x86.GetImageType(mapping.imageType)
		check(err)

		converted, err := defs.ImageTypeFromConfig(config, qcow2)
		check(err)

		basename := fmt.Sprintf("%s-%s-%s", mapping.distro, mapping.arch, mapping.imageType)
		compare(doManifest(converted), doManifest(qcow2), basename)
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
	check(err)

	archName := imgType.Arch().Name()
	depsolvedSets, err := manifestmock.Depsolve(common.Must(mf.GetPackageSetChains()), archName, nil, false)
	check(err)
	containerSpecs := manifestmock.ResolveContainers(mf.GetContainerSourceSpecs())
	commitSpecs := manifestmock.ResolveCommits(mf.GetOSTreeSourceSpecs())
	flatpakSpecs := manifestmock.ResolveFlatpaks(mf.GetFlatpakSourceSpecs())
	mfs, err := mf.Serialize(depsolvedSets, containerSpecs, commitSpecs, flatpakSpecs, nil)
	check(err)
	return mfs
}
