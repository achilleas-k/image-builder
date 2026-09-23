package defs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/osbuild/image-builder/internal/common"
	"github.com/osbuild/image-builder/pkg/arch"
	"github.com/osbuild/image-builder/pkg/datasizes"
	"github.com/osbuild/image-builder/pkg/disk"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/osbuild"
	"github.com/osbuild/image-builder/pkg/platform"
	"github.com/osbuild/image-builder/pkg/rpmmd"
	"go.yaml.in/yaml/v3"
)

type BuildConfig struct {
	Names          []string        `yaml:"names,omitempty"`
	Filename       string          `yaml:"filename,omitempty"`
	MIMEType       string          `yaml:"mime_type,omitempty"`
	Bootable       bool            `yaml:"bootable,omitempty"`
	DefaultSize    datasizes.Size  `yaml:"default_size,omitempty"`
	ImageFunc      imageFunc       `yaml:"image_func,omitempty"`
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
	InstallWeakDeps bool                        `yaml:"install_weak_deps,omitempty"`
	Sets            map[string]rpmmd.PackageSet `yaml:"sets,omitempty"`
}

type SystemConfig struct {
	KernelOptions          []string         `yaml:"kernel_options,omitempty"`
	DefaultOSCAPDatastream *string          `yaml:"default_oscap_datastream,omitempty"`
	Locale                 string           `yaml:"locale,omitempty"`
	DefaultKernel          string           `yaml:"default_kernel,omitempty"`
	UpdateDefaultKernel    bool             `yaml:"update_default_kernel,omitempty"`
	Hostname               string           `yaml:"hostname,omitempty"`
	Systemd                SystemdConfig    `yaml:"systemd,omitempty"`
	Init                   InitConfig       `yaml:"init,omitempty"`
	Keyboard               KeyboardConfig   `yaml:"keyboard,omitempty"`
	Time                   TimeConfig       `yaml:"time,omitempty"`
	Modprobe               ModprobeConfigs  `yaml:"modprobe,omitempty"`
	Network                NetworkConfig    `yaml:"network,omitempty"`
	CloudInit              CloudInitConfigs `yaml:"cloud_init,omitempty"`
	Sshd                   SshdConfig       `yaml:"sshd,omitempty"`
}

type PartitionTables struct {
	RequiredPartitionSizes map[string]datasizes.Size      `yaml:"required_partition_sizes,omitempty"`
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
	DefaultTarget          string          `yaml:"default_target,omitempty"`
	MachineIdUninitialized bool            `yaml:"machine_id_uninitialized,omitempty"`
	Services               SystemdServices `yaml:"services,omitempty"`
	Logind                 SystemdLogind   `yaml:"logind,omitempty"`
	Dropins                SystemdDropins  `yaml:"dropins,omitempty"`
}

type SystemdServices struct {
	Enable  []string `yaml:"enable,omitempty"`
	Disable []string `yaml:"disable,omitempty"`
	Mask    []string `yaml:"mask,omitempty"`
}

type SystemdLogind struct {
	Configs SystemdLogindDropins `yaml:"configs,omitempty"`
}

type SystemdLogindDropins []SystemdLogindDropin

func (dropins SystemdLogindDropins) stageOptions() []*osbuild.SystemdLogindStageOptions {
	if len(dropins) == 0 {
		return nil
	}

	options := make([]*osbuild.SystemdLogindStageOptions, len(dropins))
	for idx, dropin := range dropins {
		options[idx] = &osbuild.SystemdLogindStageOptions{
			Filename: dropin.Filename,
			Config: osbuild.SystemdLogindConfigDropin{
				Login: osbuild.SystemdLogindConfigLoginSection{
					NAutoVTs:  dropin.Config.NAutoVTs,
					ReserveVT: dropin.Config.ReserveVT,
				},
			},
		}
	}
	return options
}

type SystemdLogindDropin struct {
	Filename string       `yaml:"filename,omitempty"`
	Config   LogindConfig `yaml:"config,omitempty"`
}

type LogindConfig struct {
	NAutoVTs  *int `yaml:"NAutoVTs,omitempty"`
	ReserveVT *int `yaml:"ReserveVT,omitempty"`
}

type InitConfig struct {
	Dracut Dracuts `yaml:"dracut,omitempty"`
}

type Dracuts []Dracut

type Dracut struct {
	Filename string       `yaml:"filename,omitempty"`
	Config   DracutConfig `yaml:"config,omitempty"`
}

type DracutConfig struct {
	AddDrivers []string `yaml:"add_drivers,omitempty"`
}

func (dracuts Dracuts) stageOptions() []*osbuild.DracutConfStageOptions {
	options := make([]*osbuild.DracutConfStageOptions, len(dracuts))
	for idx, dracut := range dracuts {
		options[idx] = &osbuild.DracutConfStageOptions{
			Filename: dracut.Filename,
			Config: osbuild.DracutConfigFile{
				Compress:       "",
				Modules:        []string{},
				AddModules:     []string{},
				OmitModules:    []string{},
				Drivers:        []string{},
				AddDrivers:     dracut.Config.AddDrivers,
				ForceDrivers:   []string{},
				Filesystems:    []string{},
				Install:        []string{},
				EarlyMicrocode: nil,
				Reproducible:   nil,
			},
		}
	}

	return options
}

type KeyboardConfig struct {
	Keymap  string   `yaml:"keymap,omitempty"`
	Layouts []string `yaml:"layouts,omitempty"`
}

func (keyboard KeyboardConfig) stageOptions() *osbuild.KeymapStageOptions {
	if len(keyboard.Layouts) == 0 && keyboard.Keymap == "" {
		return nil
	}
	options := &osbuild.KeymapStageOptions{
		Keymap: keyboard.Keymap,
	}

	if len(keyboard.Layouts) > 0 {
		options.X11Keymap = &osbuild.X11KeymapOptions{
			Layouts: keyboard.Layouts,
		}
	}

	return options
}

type TimeConfig struct {
	Timezone   string     `yaml:"timezone,omitempty"`
	NTPServers NTPServers `yaml:"ntp_servers,omitempty"`
	LeapsecTz  *string    `yaml:"leapsectz,omitempty"`
}

func (time TimeConfig) chronyStageOptions() *osbuild.ChronyStageOptions {
	if len(time.NTPServers) == 0 && time.LeapsecTz == nil {
		return nil
	}

	return &osbuild.ChronyStageOptions{
		Servers:   time.NTPServers.stageOptions(),
		LeapsecTz: time.LeapsecTz,
	}
}

type NTPServers []NTPServer

type NTPServer struct {
	Hostname string `yaml:"hostname,omitempty"`
	Minpoll  *int   `yaml:"minpoll,omitempty"`
	Maxpoll  *int   `yaml:"maxpoll,omitempty"`
	Iburst   *bool  `yaml:"iburst,omitempty"`
	Prefer   *bool  `yaml:"prefer,omitempty"`
}

type ModprobeConfigs []ModprobeConfig

func (configs ModprobeConfigs) stageOptions() []*osbuild.ModprobeStageOptions {
	options := make([]*osbuild.ModprobeStageOptions, len(configs))
	for idx, config := range configs {
		options[idx] = config.stageOptions()
	}
	return options
}

type ModprobeConfig struct {
	Filename string           `yaml:"filename,omitempty"`
	Commands ModprobeCommands `yaml:"commands,omitempty"`
}

func (config ModprobeConfig) stageOptions() *osbuild.ModprobeStageOptions {
	return &osbuild.ModprobeStageOptions{
		Filename: config.Filename,
		Commands: config.Commands.stageOptions(),
	}
}

type ModprobeCommands []ModprobeCommand

func (commands ModprobeCommands) stageOptions() osbuild.ModprobeConfigCmdList {
	options := make([]osbuild.ModprobeConfigCmd, len(commands))

	for idx, command := range commands {
		var cmdOptions osbuild.ModprobeConfigCmd
		switch command.Command {
		case "blacklist":
			cmdOptions = osbuild.NewModprobeConfigCmdBlacklist(command.ModuleName)
		case "install":
			cmdOptions = osbuild.NewModprobeConfigCmdInstall(command.ModuleName, command.Cmdline)
		}
		options[idx] = cmdOptions
	}

	return options
}

type ModprobeCommand struct {
	Command    string `yaml:"command,omitempty"`
	ModuleName string `yaml:"modulename,omitempty"`
	Cmdline    string `yaml:"cmdline,omitempty"`
}

func (servers NTPServers) stageOptions() []osbuild.ChronyConfigServer {
	if len(servers) == 0 {
		return nil
	}
	chronyServers := make([]osbuild.ChronyConfigServer, len(servers))
	for idx, server := range servers {
		chronyServers[idx] = osbuild.ChronyConfigServer{
			Hostname: server.Hostname,
			Minpoll:  server.Minpoll,
			Maxpoll:  server.Maxpoll,
			Iburst:   server.Iburst,
			Prefer:   server.Prefer,
		}
	}
	return chronyServers
}

type NetworkConfig struct {
	Enabled    bool `yaml:"enabled,omitempty"`
	NoZeroConf bool `yaml:"no_zero_conf,omitempty"`
}

type CloudInitConfigs []CloudInit

func (cfgs CloudInitConfigs) stageOptions() []*osbuild.CloudInitStageOptions {
	if len(cfgs) == 0 {
		return nil
	}

	options := make([]*osbuild.CloudInitStageOptions, len(cfgs))
	for idx, cfg := range cfgs {
		options[idx] = cfg.stageOptions()
	}

	return options
}

type CloudInit struct {
	Filename string          `yaml:"filename,omitempty"`
	Config   CloudInitConfig `yaml:"config,omitempty"`
}

func (cfg CloudInit) stageOptions() *osbuild.CloudInitStageOptions {
	return &osbuild.CloudInitStageOptions{
		Filename: cfg.Filename,
		Config:   cfg.Config.stageOptions(),
	}
}

type CloudInitConfig struct {
	SystemInfo     *CloudInitConfigSystemInfo `yaml:"system_info,omitempty"`
	DatasourceList []string                   `yaml:"datasource_list,omitempty"`
}

func (cfg CloudInitConfig) stageOptions() osbuild.CloudInitConfigFile {
	return osbuild.CloudInitConfigFile{
		SystemInfo: &osbuild.CloudInitConfigSystemInfo{
			DefaultUser: &osbuild.CloudInitConfigDefaultUser{
				Name: cfg.SystemInfo.DefaultUser.Name,
			},
		},
	}
}

type CloudInitConfigSystemInfo struct {
	Name        string                      `yaml:"name,omitempty"`
	DefaultUser *CloudInitConfigDefaultUser `yaml:"default_user,omitempty"`
}

type CloudInitConfigDefaultUser struct {
	Name string `yaml:"name,omitempty"`
}

type SshdConfig struct {
	PasswordAuthentication *bool `yaml:"password_authentication,omitempty"`
}

func (cfg SshdConfig) stageOptions() *osbuild.SshdConfigStageOptions {
	if cfg.PasswordAuthentication == nil {
		return nil
	}

	return &osbuild.SshdConfigStageOptions{
		Config: osbuild.SshdConfigConfig{
			PasswordAuthentication: cfg.PasswordAuthentication,
		},
	}
}

type SystemdDropins []SystemdDropin

func (dropins SystemdDropins) stageOptions() []*osbuild.SystemdUnitStageOptions {
	if len(dropins) == 0 {
		return nil
	}

	options := make([]*osbuild.SystemdUnitStageOptions, len(dropins))
	for idx, dropin := range dropins {
		options[idx] = dropin.stageOptions()
	}
	return options
}

type SystemdDropin struct {
	Unit   string              `yaml:"unit,omitempty"`
	Dropin string              `yaml:"dropin,omitempty"`
	Config SystemdDropinConfig `yaml:"config,omitempty"`
}

func (dropin SystemdDropin) stageOptions() *osbuild.SystemdUnitStageOptions {
	return &osbuild.SystemdUnitStageOptions{
		Unit:   dropin.Unit,
		Dropin: dropin.Dropin,
		Config: dropin.Config.stageOptions(),
	}
}

type SystemdDropinConfig struct {
	Service SystemdDropinService `yaml:"service,omitempty"`
}

func (cfg SystemdDropinConfig) stageOptions() osbuild.SystemdServiceUnitDropin {
	var vars []osbuild.EnvironmentVariable
	if len(cfg.Service.Environment) > 0 {
		vars = make([]osbuild.EnvironmentVariable, len(cfg.Service.Environment))
		for idx, v := range cfg.Service.Environment {
			vars[idx] = osbuild.EnvironmentVariable(v)
		}
	}
	return osbuild.SystemdServiceUnitDropin{
		Service: &osbuild.SystemdUnitServiceSection{
			Environment: vars,
		},
	}
}

type SystemdDropinService struct {
	Environment []KeyValue `yaml:"environment,omitempty"`
}

type KeyValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
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

// ImageTypeFromConfig is a temporary conversion function for verifying image configuration parity.
func ImageTypeFromConfig(config BuildConfig, borrowFrom distro.ImageType) (distro.ImageType, error) {
	bf := borrowFrom.(*imageType)

	pt := disk.PartitionTable(config.PartitionTable.Architecture[config.Architecture])

	packageSets := make(map[string]rpmmd.PackageSet, len(config.Content.Packages.Sets))
	for name, set := range config.Content.Packages.Sets {
		packageSets[name] = rpmmd.PackageSet{
			Include: slices.Clone(set.Include),
			Exclude: slices.Clone(set.Exclude),
		}
	}

	// NOTE .stageOptions() methods are temporary, for serving this conversion only
	it := imageType{
		name:                   config.Names[0],
		nameAliases:            config.Names[1:],
		arch:                   bf.arch,
		filename:               config.Filename,
		mimeType:               config.MIMEType,
		bootable:               config.Bootable,
		defaultSize:            config.DefaultSize,
		exports:                config.Exports,
		requiredPartitionSizes: config.PartitionTable.RequiredPartitionSizes,
		partitionTable:         &pt,
		imageConfig: distro.ImageConfig{
			KernelOptions:          config.System.KernelOptions,
			Locale:                 &config.System.Locale,
			Hostname:               &config.System.Hostname,
			UpdateDefaultKernel:    &config.System.UpdateDefaultKernel,
			DefaultKernel:          &config.System.DefaultKernel,
			EnabledServices:        config.System.Systemd.Services.Enable,
			DisabledServices:       config.System.Systemd.Services.Disable,
			MaskedServices:         config.System.Systemd.Services.Mask,
			DefaultTarget:          &config.System.Systemd.DefaultTarget,
			MachineIdUninitialized: &config.System.Systemd.MachineIdUninitialized,
			InstallWeakDeps:        &config.Content.Packages.InstallWeakDeps,
			DracutConf:             config.System.Init.Dracut.stageOptions(),
			Keyboard:               config.System.Keyboard.stageOptions(),
			Timezone:               &config.System.Time.Timezone,
			TimeSynchronization:    config.System.Time.chronyStageOptions(),
			Modprobe:               config.System.Modprobe.stageOptions(),
			Sysconfig: &distro.Sysconfig{
				Networking: config.System.Network.Enabled,
				NoZeroConf: config.System.Network.NoZeroConf,
			},
			SystemdLogind: config.System.Systemd.Logind.Configs.stageOptions(),
			CloudInit:     config.System.CloudInit.stageOptions(),
			SshdConfig:    config.System.Sshd.stageOptions(),
			SystemdDropin: config.System.Systemd.Dropins.stageOptions(),
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

// stupid wrapper to make calls less visible and focus on the mapping
func v[T any](p *T) T {
	return common.ValueOrEmpty(p)
}

func NewConfig(it *imageType) (BuildConfig, error) {

	bc := BuildConfig{
		Names:       append([]string{it.name}, it.nameAliases...),
		Filename:    it.filename,
		MIMEType:    it.mimeType,
		Bootable:    it.bootable,
		DefaultSize: it.defaultSize,
		ImageFunc:   it.image,
		Exports:     it.exports,
		Blueprint: BlueprintConfig{
			SupportedOptions: it.blueprint.SupportedOptions,
		},
		Content: ContentConfig{
			Packages: PackagesConfig{
				InstallWeakDeps: v(it.installWeakDeps),
				Sets:            it.packageSets,
			},
		},
		System: SystemConfig{
			KernelOptions:          it.imageConfig.KernelOptions,
			DefaultOSCAPDatastream: it.imageConfig.DefaultOSCAPDatastream,
			Locale:                 v(it.imageConfig.Locale),
			Time: TimeConfig{
				Timezone: v(it.imageConfig.Timezone),
			},
			DefaultKernel:       v(it.imageConfig.DefaultKernel),
			UpdateDefaultKernel: v(it.imageConfig.UpdateDefaultKernel),
			Hostname:            v(it.imageConfig.Hostname),
			Systemd: SystemdConfig{
				DefaultTarget:          v(it.imageConfig.DefaultTarget),
				MachineIdUninitialized: v(it.imageConfig.MachineIdUninitialized),
				Services: SystemdServices{
					Enable:  it.imageConfig.EnabledServices,
					Disable: it.imageConfig.DisabledServices,
					Mask:    it.imageConfig.MaskedServices,
				},
			},
		},
		PartitionTable: PartitionTables{
			RequiredPartitionSizes: it.requiredPartitionSizes,
			Architecture: map[string]disk.PartitionTable{
				it.arch.Name(): v(it.partitionTable),
			},
		},
		DistributionName: it.arch.distro.Name(),
		Depsolver:        "osbuild-depsolve-dnf",
		PackageManager:   "rpm",
		UEFIVendor:       it.platform.GetUEFIVendor(),
		Architecture:     it.arch.Name(),
		Bootloader:       it.platform.GetBootloader().String(),
	}

	return bc, nil
}

// PLATFORM INTERFACE
func (bc BuildConfig) GetArch() arch.Arch {
	a, err := arch.FromString(bc.Architecture)
	common.PanicOnError(err)
	return a
}

func (bc BuildConfig) GetImageFormat() platform.ImageFormat {
	// TODO: add to BuildConfig
	return platform.FORMAT_UNSET
}

func (bc BuildConfig) GetQCOW2Compat() string {
	// TODO: add to BuildConfig
	return ""
}

func (bc BuildConfig) GetBIOSPlatform() string {
	// TODO: add to BuildConfig
	return ""
}

func (bc BuildConfig) GetUEFIVendor() string {
	return bc.UEFIVendor
}

func (bc BuildConfig) GetExtraUEFIArchitectures() []string {
	// TODO: add to BuildConfig
	return nil
}

func (bc BuildConfig) GetZiplSupport() bool {
	return false
}

func (bc BuildConfig) GetPackages() []string {
	return nil
}

func (bc BuildConfig) GetBuildPackages() []string {
	return nil
}

func (bc BuildConfig) GetBootFiles() []platform.BootFile {
	return nil
}
func (bc BuildConfig) GetBootloader() platform.Bootloader {
	bootloader, err := platform.FromString(bc.Bootloader)
	common.PanicOnError(err)
	return bootloader
}

func (bc BuildConfig) GetFIPSMenu() bool {
	// TODO: add to BuildConfig
	return false
}

// ENVIRONMENT INTERFACE
func (bc BuildConfig) GetRepos() []rpmmd.RepoConfig {
	return nil
}

func (bc BuildConfig) GetServices() []string {
	return nil
}
