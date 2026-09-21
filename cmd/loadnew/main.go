package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/osbuild/image-builder/pkg/distro/defs"
)

func jsonPrint(thingie any) {
	c, err := json.MarshalIndent(thingie, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(c))
}

func main() {
	paths := []string{
		"distro/fedora.yaml",
		"platform/x86_64-base.yaml",
		"type/qcow2.yaml",
	}

	config, err := defs.Load(paths...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	jsonPrint(config)

	fedora, err := defs.New("fedora-44")
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	x86, err := fedora.GetArch("x86_64")
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	qcow2, err := x86.GetImageType("generic-qcow2")
	if err != nil {
		fmt.Fprintf(os.Stderr, "E: %s\n", err)
		os.Exit(1)
	}

	jsonPrint(qcow2)
}
