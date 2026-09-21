package defs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/osbuild/image-builder/pkg/arch"
	"github.com/osbuild/image-builder/pkg/datasizes"
	"github.com/osbuild/image-builder/pkg/distro"
	"go.yaml.in/yaml/v3"
)

type BuildConfig struct {
	Names          []string        `yaml:"names,omitempty"`
	Filename       string          `yaml:"filename,omitempty"`
	MIMEType       string          `yaml:"mime_type,omitempty"`
	Bootable       bool            `yaml:"bootable,omitempty"`
	DefaultSize    string          `yaml:"default_size,omitempty"`
	ImageFunc      string          `yaml:"image_func,omitempty"`
	Exports        []string        `yaml:"exports,omitempty"`
	Blueprint      BlueprintConfig `yaml:"blueprint,omitempty"`
	Content        ContentConfig   `yaml:"content,omitempty"`
	Config         SystemConfig    `yaml:"config,omitempty"`
	PartitionTable PartitionTables `yaml:"partition_table,omitempty"`

	DistributionName string `yaml:"distribution_name"`
	Depsolver        string `yaml:"depsolver"`
	PackageManager   string `yaml:"package_manager"`
	UEFIVendor       string `yaml:"uefi_vendor"`

	Architecture string `yaml:"architecture"`
	Bootloader   string `yaml:"bootloader"`
}

type BlueprintConfig struct {
	SupportedOptions []string `yaml:"supported_options,omitempty"`
}

type ContentConfig struct {
	Packages map[string]PackageSet `yaml:"packages,omitempty"`
}

type PackageSet struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type SystemConfig struct {
	DefaultTarget string   `yaml:"default_target,omitempty"`
	KernelOptions []string `yaml:"kernel_options,omitempty"`
}

type PartitionTables struct {
	RequiredPartitionSizes map[string]string               `yaml:"required_partition_sizes,omitempty"`
	Architecture           map[string]PartitionTableConfig `yaml:"architecture,omitempty"`
}

type PartitionTableConfig struct {
	UUID       string      `yaml:"uuid,omitempty"`
	Type       string      `yaml:"type,omitempty"`
	Partitions []Partition `yaml:"partitions,omitempty"`
}

type Partition struct {
	Size        string         `yaml:"size,omitempty"`
	Bootable    bool           `yaml:"bootable,omitempty"`
	Type        string         `yaml:"type,omitempty"`
	PayloadType string         `yaml:"payload_type,omitempty"`
	Payload     *PayloadConfig `yaml:"payload,omitempty"`
}

type PayloadConfig struct {
	Type         string `yaml:"type,omitempty"`
	Mountpoint   string `yaml:"mountpoint,omitempty"`
	Label        string `yaml:"label,omitempty"`
	FstabOptions string `yaml:"fstab_options,omitempty"`
	FstabFreq    int    `yaml:"fstab_freq"`
	FstabPassno  int    `yaml:"fstab_passno"`
}

func Load(paths ...string) (BuildConfig, error) {
	var config BuildConfig
	for _, path := range paths {
		fp, err := os.Open(filepath.Join("data/configs", path))
		if err != nil {
			return BuildConfig{}, err
		}
		decoder := yaml.NewDecoder(fp)
		decoder.KnownFields(true)

		if err := decoder.Decode(&config); err != nil {
			return BuildConfig{}, fmt.Errorf("%s: %w", path, err)
		}
	}

	return config, nil
}

func ImageTypeFromConfig(config BuildConfig) (distro.ImageType, error) {
	a, err := arch.FromString(config.Architecture)
	if err != nil {
		return nil, err
	}
	arch := &architecture{
		arch: a,
	}

	var size *datasizes.Size
	if err := size.UnmarshalText([]byte(config.DefaultSize)); err != nil {
		return nil, err
	}

	partSizes := make(map[string]datasizes.Size, len(config.PartitionTable.RequiredPartitionSizes))
	for path, sizeStr := range config.PartitionTable.RequiredPartitionSizes {
		partSizes[path] = convertSize(sizeStr)
	}

	it := imageType{
		name:                   config.Names[0],
		nameAliases:            config.Names[1:],
		arch:                   arch,
		filename:               config.Filename,
		mimeType:               config.MIMEType,
		bootable:               config.Bootable,
		defaultSize:            convertSize(config.DefaultSize),
		exports:                config.Exports,
		requiredPartitionSizes: partSizes,

		blueprint: blueprintOptions{
			// The blueprint contains a few fields that are essentially
			// metadata and not configuration / customizations. These should
			// always be implicitly supported by all image types.
			SupportedOptions: append(slices.Clone(config.Blueprint.SupportedOptions), "name", "version", "description"),
		},
	}

	return &it, nil
}

func convertSize(s string) datasizes.Size {
	var size *datasizes.Size
	if err := size.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return *size
}
