package defs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/osbuild/image-builder/internal/common"
	"github.com/osbuild/image-builder/pkg/datasizes"
	"github.com/osbuild/image-builder/pkg/disk"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/rpmmd"
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
	System         SystemConfig    `yaml:"system,omitempty"`
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
	Packages PackagesConfig `yaml:"packages,omitempty"`
}

type PackagesConfig struct {
	InstallWeakDeps bool                  `yaml:"install_weak_deps,omitempty"`
	Sets            map[string]PackageSet `yaml:"sets,omitempty"`
}

type PackageSet struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type SystemConfig struct {
	KernelOptions          []string      `yaml:"kernel_options,omitempty"`
	DefaultOSCAPDatastream *string       `yaml:"default_oscap_datastream,omitempty"`
	Locale                 string        `yaml:"locale,omitempty"`
	Timezone               string        `yaml:"timezone,omitempty"`
	DefaultKernel          string        `yaml:"default_kernel,omitempty"`
	UpdateDefaultKernel    bool          `yaml:"update_default_kernel,omitempty"`
	Hostname               string        `yaml:"hostname,omitempty"`
	Systemd                SystemdConfig `yaml:"systemd,omitempty"`
}

type PartitionTables struct {
	RequiredPartitionSizes map[string]string              `yaml:"required_partition_sizes,omitempty"`
	Architecture           map[string]disk.PartitionTable `yaml:"architecture,omitempty"`
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

type SystemdConfig struct {
	DefaultTarget          string                `yaml:"default_target,omitempty"`
	MachineIdUninitialized bool                  `yaml:"machine_id_uninitialized,omitempty"`
	Services               SystemdServicesConfig `yaml:"services,omitempty"`
}

type SystemdServicesConfig struct {
	Enable  []string `yaml:"enable,omitempty"`
	Disable []string `yaml:"disable,omitempty"`
	Mask    []string `yaml:"mask,omitempty"`
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

func ImageTypeFromConfig(config BuildConfig, borrowFrom distro.ImageType) (distro.ImageType, error) {
	bf := borrowFrom.(*imageType)
	partSizes := make(map[string]datasizes.Size, len(config.PartitionTable.RequiredPartitionSizes))
	for path, sizeStr := range config.PartitionTable.RequiredPartitionSizes {
		partSizes[path] = convertSize(sizeStr)
	}

	pt := disk.PartitionTable(config.PartitionTable.Architecture[config.Architecture])

	packageSets := make(map[string]rpmmd.PackageSet, len(config.Content.Packages.Sets))
	for name, set := range config.Content.Packages.Sets {
		packageSets[name] = rpmmd.PackageSet{
			Include: slices.Clone(set.Include),
			Exclude: slices.Clone(set.Exclude),
		}
	}

	it := imageType{
		name:                   config.Names[0],
		nameAliases:            config.Names[1:],
		arch:                   bf.arch,
		filename:               config.Filename,
		mimeType:               config.MIMEType,
		bootable:               config.Bootable,
		defaultSize:            convertSize(config.DefaultSize),
		exports:                config.Exports,
		requiredPartitionSizes: partSizes,
		partitionTable:         &pt,
		imageConfig: distro.ImageConfig{
			KernelOptions:          config.System.KernelOptions,
			Locale:                 &config.System.Locale,
			Hostname:               &config.System.Hostname,
			Timezone:               &config.System.Timezone,
			UpdateDefaultKernel:    common.ToPtr(true),
			DefaultKernel:          common.ToPtr("kernel-core"),
			EnabledServices:        config.System.Systemd.Services.Enable,
			DisabledServices:       config.System.Systemd.Services.Disable,
			MaskedServices:         config.System.Systemd.Services.Mask,
			DefaultTarget:          &config.System.Systemd.DefaultTarget,
			MachineIdUninitialized: &config.System.Systemd.MachineIdUninitialized,
			InstallWeakDeps:        &config.Content.Packages.InstallWeakDeps,
		},

		platform: bf.platform,

		blueprint: blueprintOptions{
			// The blueprint contains a few fields that are essentially
			// metadata and not configuration / customizations. These should
			// always be implicitly supported by all image types.
			SupportedOptions: append(slices.Clone(config.Blueprint.SupportedOptions), "name", "version", "description"),
		},
		image: diskImage,

		packageSets: packageSets,
	}

	return &it, nil
}

func convertSize(s string) datasizes.Size {
	size := new(datasizes.Size)
	if err := size.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return *size
}
