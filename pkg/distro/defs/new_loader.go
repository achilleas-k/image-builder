package defs

import (
	"math/rand"

	"github.com/osbuild/blueprint/pkg/blueprint"
	"github.com/osbuild/image-builder/internal/environment"
	"github.com/osbuild/image-builder/pkg/container"
	"github.com/osbuild/image-builder/pkg/customizations/fsnode"
	"github.com/osbuild/image-builder/pkg/customizations/oci"
	"github.com/osbuild/image-builder/pkg/customizations/ostreeserver"
	"github.com/osbuild/image-builder/pkg/customizations/shell"
	"github.com/osbuild/image-builder/pkg/customizations/subscription"
	"github.com/osbuild/image-builder/pkg/customizations/users"
	"github.com/osbuild/image-builder/pkg/customizations/wsl"
	"github.com/osbuild/image-builder/pkg/datasizes"
	"github.com/osbuild/image-builder/pkg/disk"
	"github.com/osbuild/image-builder/pkg/disk/partition"
	"github.com/osbuild/image-builder/pkg/distro"
	"github.com/osbuild/image-builder/pkg/image"
	"github.com/osbuild/image-builder/pkg/manifest"
	"github.com/osbuild/image-builder/pkg/osbuild"
	"github.com/osbuild/image-builder/pkg/rpmmd"
)

func loadImageConfiguration() (*imageType, error) {
	it := imageType{
		name:        "",
		nameAliases: []string{},
		filename:    "",
		mimeType:    "",
		compression: "",
		packageSets: map[string]rpmmd.PackageSet{
			"": {
				Include:         []string{},
				Exclude:         []string{},
				EnabledModules:  []string{},
				Repositories:    []rpmmd.RepoConfig{},
				InstallWeakDeps: false,
			},
		},
		partitionTable: &disk.PartitionTable{
			Size:                0,
			UUID:                "",
			Type:                0,
			Partitions:          []disk.Partition{},
			GrainSize:           0,
			SectorSize:          0,
			ExtraPadding:        0,
			StartOffset:         0,
			AbsoluteStartOffset: false,
			Policy: &disk.PartitionTablePolicy{
				EnsureXBOOTLDR:     false,
				GrowRootToFillDisk: nil,
			},
		},
		imageConfig: distro.ImageConfig{
			Hostname: nil,
			Timezone: nil,
			TimeSynchronization: &osbuild.ChronyStageOptions{
				Servers:   []osbuild.ChronyConfigServer{},
				Refclocks: []osbuild.ChronyConfigRefclock{},
				LeapsecTz: nil,
			},
			Locale: nil,
			Keyboard: &osbuild.KeymapStageOptions{
				Keymap: "",
				X11Keymap: &osbuild.X11KeymapOptions{
					Layouts: []string{},
				},
			},
			EnabledServices:  []string{},
			DisabledServices: []string{},
			MaskedServices:   []string{},
			DefaultTarget:    nil,
			Sysconfig: &distro.Sysconfig{
				Networking:                  false,
				NoZeroConf:                  false,
				CreateDefaultNetworkScripts: false,
			},
			DefaultKernel:       nil,
			UpdateDefaultKernel: nil,
			KernelOptions:       []string{},
			DefaultKernelName:   nil,
			GPGKeyFiles:         []string{},
			NoSELinux:           nil,
			SELinuxForceRelabel: nil,
			ExcludeDocs:         nil,
			ShellInit:           []shell.InitFile{},
			RHSMConfig: map[subscription.RHSMStatus]*subscription.RHSMConfig{
				"": {
					DnfPlugins: subscription.SubManDNFPluginsConfig{
						ProductID: subscription.DNFPluginConfig{
							Enabled: nil,
						},
						SubscriptionManager: subscription.DNFPluginConfig{
							Enabled: nil,
						},
					},
					YumPlugins: subscription.SubManDNFPluginsConfig{
						ProductID: subscription.DNFPluginConfig{
							Enabled: nil,
						},
						SubscriptionManager: subscription.DNFPluginConfig{
							Enabled: nil,
						},
					},
					SubMan: subscription.SubManConfig{
						Rhsm: subscription.SubManRHSMConfig{
							ManageRepos:          nil,
							AutoEnableYumPlugins: nil,
						},
						Rhsmcertd: subscription.SubManRHSMCertdConfig{
							AutoRegistration: nil,
						},
					},
				},
			},
			SystemdLogind: []*osbuild.SystemdLogindStageOptions{},
			CloudInit:     []*osbuild.CloudInitStageOptions{},
			Modprobe:      []*osbuild.ModprobeStageOptions{},
			DracutConf:    []*osbuild.DracutConfStageOptions{},
			SystemdDropin: []*osbuild.SystemdUnitStageOptions{},
			SystemdUnit:   []*osbuild.SystemdUnitCreateStageOptions{},
			Authselect: &osbuild.AuthselectStageOptions{
				Profile:  "",
				Features: []string{},
			},
			SELinuxConfig: &osbuild.SELinuxConfigStageOptions{
				State: "",
				Type:  "",
			},
			Tuned: &osbuild.TunedStageOptions{
				Profiles: []string{},
			},
			Tmpfilesd:     []*osbuild.TmpfilesdStageOptions{},
			PamLimitsConf: []*osbuild.PamLimitsConfStageOptions{},
			Sysctld:       []*osbuild.SysctldStageOptions{},
			DNFConfig: &distro.DNFConfig{
				Options: &osbuild.DNFConfigStageOptions{
					Variables: []osbuild.DNFVariable{},
					Config: &osbuild.DNFConfig{
						Main: &osbuild.DNFConfigMain{
							IPResolve: "",
						},
					},
				},
				SetReleaseverVar: nil,
			},
			SshdConfig: &osbuild.SshdConfigStageOptions{
				Config: osbuild.SshdConfigConfig{
					PasswordAuthentication:          nil,
					ChallengeResponseAuthentication: nil,
					ClientAliveInterval:             nil,
					PermitRootLogin:                 nil,
				},
			},
			Authconfig: &osbuild.AuthconfigStageOptions{},
			PwQuality: &osbuild.PwqualityConfStageOptions{
				Config: osbuild.PwqualityConfConfig{
					Minlen:   nil,
					Dcredit:  nil,
					Ucredit:  nil,
					Lcredit:  nil,
					Ocredit:  nil,
					Minclass: nil,
				},
			},
			WAAgentConfig: &osbuild.WAAgentConfStageOptions{
				Config: osbuild.WAAgentConfig{
					ProvisioningUseCloudInit: nil,
					ProvisioningEnabled:      nil,
					RDFormat:                 nil,
					RDEnableSwap:             nil,
				},
			},
			Grub2Config: &osbuild.GRUB2Config{
				Default:         "",
				DisableRecovery: nil,
				DisableSubmenu:  nil,
				Distributor:     "",
				Terminal:        []string{},
				TerminalInput:   []string{},
				TerminalOutput:  []string{},
				Timeout:         0,
				TimeoutStyle:    "",
				Serial:          "",
			},
			DNFAutomaticConfig: &osbuild.DNFAutomaticConfigStageOptions{
				Config: &osbuild.DNFAutomaticConfig{
					Commands: &osbuild.DNFAutomaticConfigCommands{
						ApplyUpdates: nil,
						UpgradeType:  "",
					},
				},
			},
			YumConfig: &osbuild.YumConfigStageOptions{
				Config: &osbuild.YumConfigConfig{
					HttpCaching: nil,
				},
				Plugins: &osbuild.YumConfigPlugins{
					Langpacks: &osbuild.YumConfigPluginsLangpacks{
						Locales: []string{},
					},
				},
			},
			YUMRepos: []*osbuild.YumReposStageOptions{},
			Firewall: &osbuild.FirewallStageOptions{
				Ports:            []string{},
				EnabledServices:  []string{},
				DisabledServices: []string{},
				DefaultZone:      "",
				Zones:            []osbuild.FirewallZone{},
			},
			UdevRules: &osbuild.UdevRulesStageOptions{
				Filename: "",
				Rules:    []osbuild.UdevRule{},
			},
			GCPGuestAgentConfig: &osbuild.GcpGuestAgentConfigOptions{
				ConfigScope: "",
				Config: &osbuild.GcpGuestAgentConfig{
					Accounts: &osbuild.GcpGuestAgentConfigAccounts{
						DeprovisionRemove: nil,
						Groups:            []string{},
						UseraddCmd:        "",
						UserdelCmd:        "",
						UsermodCmd:        "",
						GpasswdAddCmd:     "",
						GpasswdRemoveCmd:  "",
						GroupaddCmd:       "",
					},
					Daemons: &osbuild.GcpGuestAgentConfigDaemons{
						AccountsDaemon:  nil,
						ClockSkewDaemon: nil,
						NetworkDaemon:   nil,
					},
					InstanceSetup: &osbuild.GcpGuestAgentConfigInstanceSetup{
						HostKeyTypes:     []string{},
						OptimizeLocalSsd: nil,
						NetworkEnabled:   nil,
						SetBotoConfig:    nil,
						SetHostKeys:      nil,
						SetMultiqueue:    nil,
					},
					IpForwarding: &osbuild.GcpGuestAgentConfigIpForwarding{
						EthernetProtoId:   "",
						IpAliases:         nil,
						TargetInstanceIps: nil,
					},
					MetadataScripts: &osbuild.GcpGuestAgentConfigMetadataScripts{
						DefaultShell: "",
						RunDir:       "",
						Startup:      nil,
						Shutdown:     nil,
					},
					NetworkInterfaces: &osbuild.GcpGuestAgentConfigNetworkInterfaces{
						Setup:        nil,
						IpForwarding: nil,
						DhcpCommand:  "",
					},
				},
			},
			NetworkManager: &osbuild.NMConfStageOptions{
				Path: "",
				Settings: osbuild.NMConfStageSettings{
					Main: &osbuild.NMConfSettingsMain{
						NoAutoDefault: []string{},
						Plugins:       []string{},
					},
					Device:          []osbuild.NMConfSettingsDevice{},
					GlobalDNSDomain: []osbuild.NMConfSettingsGlobalDNSDomain{},
					Keyfile: &osbuild.NMConfSettingsKeyfile{
						UnmanagedDevices: []string{},
					},
				},
			},
			Presets: []osbuild.Preset{},
			WSL: &wsl.WSL{
				Config: &wsl.WSLConfig{
					BootSystemd: false,
				},
				DistributionConfig: &wsl.WSLDistributionConfig{
					OOBE: &wsl.WSLDistributionOOBEConfig{
						DefaultUID:  nil,
						DefaultName: "",
					},
					Shortcut: &wsl.WSLDistributionShortcutConfig{
						Enabled: false,
						Icon:    "",
					},
				},
			},
			OCI: &oci.OCI{
				Archive: &oci.OCIArchiveConfig{
					Cmd:          []string{},
					Env:          []string{},
					ExposedPorts: []string{},
					User:         "",
					Labels: map[string]string{
						"": "",
					},
					StopSignal: "",
					Volumes:    []string{},
					WorkingDir: "",
				},
			},
			OSTreeServer: &ostreeserver.OSTreeServer{
				Port:       "",
				ConfigPath: "",
			},
			Users:                   []users.User{},
			Files:                   []*fsnode.File{},
			Directories:             []*fsnode.Directory{},
			Hostonly:                nil,
			KernelOptionsBootloader: nil,
			DefaultOSCAPDatastream:  nil,
			NoBLS:                   nil,
			SystemdBoot: &osbuild.SystemdBootConfig{
				RandomSeed:         "",
				MakeEntryDirectory: "",
				EntryToken:         "",
			},
			OSTreeConfSysrootReadOnly: nil,
			LockRootUser:              nil,
			IgnitionPlatform:          nil,
			InstallWeakDeps:           nil,
			InstallLangs:              []string{},
			MachineIdUninitialized:    nil,
			PermissiveRHC:             nil,
			VersionlockPackages:       []string{},
			BootupdGenMetadata:        nil,
			RPM: &distro.RPMConfig{
				Macros: []osbuild.RPMMacrosStageOptions{},
			},
		},
		installerConfig: distro.InstallerConfig{
			EnabledAnacondaModules:             []string{},
			AdditionalDracutModules:            []string{},
			AdditionalDrivers:                  []string{},
			KickstartUnattendedExtraKernelOpts: []string{},
			DefaultMenu:                        nil,
			InstallWeakDeps:                    nil,
			LoraxTemplates:                     []manifest.InstallerLoraxTemplate{},
			LoraxTemplatePackage:               nil,
			LoraxLogosPackage:                  nil,
			LoraxReleasePackage:                nil,
			ISOFiles:                           [][2]string{},
			Payload: &struct {
				Location  *manifest.PayloadLocation  "yaml:\"location,omitempty\""
				Kickstart *manifest.PayloadKickstart "yaml:\"kickstart,omitempty\""
			}{
				Location:  nil,
				Kickstart: nil,
			},
			Flatpaks: []*struct {
				Registry *struct {
					RemoteName string "yaml:\"remote_name,omitempty\""
					URL        string "yaml:\"url,omitempty\""
				} "yaml:\"registry,omitempty\""
				References []string "yaml:\"references,omitempty\""
			}{},
		},
		isoConfig: distro.ISOConfig{
			BootType:       nil,
			RootfsType:     nil,
			RootfsExcludes: []string{},
			ErofsOptions: &osbuild.ErofsStageOptions{
				Filename:     "",
				Source:       "",
				ExcludePaths: []string{},
				Compression: &osbuild.ErofsCompression{
					Method: "",
					Level:  nil,
				},
				ExtendedOptions: []string{},
				ClusterSize:     nil,
			},
			Label:        nil,
			Preparer:     nil,
			Publisher:    nil,
			Application:  nil,
			ExcludePaths: []string{},
		},
		diskConfig: distro.DiskConfig{
			MountConfiguration: nil,
			PartitioningTool:   nil,
			OVFVMWare: &osbuild.OVFVMWareStageOptions{
				OSType:                 "",
				VirtualHardwareVersion: "",
			},
		},
		environment: environment.EnvironmentConf{
			Packages: []string{},
			Repos:    []rpmmd.RepoConfig{},
			Services: []string{},
		},
		bootable:                false,
		bootISO:                 false,
		useLegacyAnacondaConfig: false,
		isoLabel:                "",
		variant:                 "",
		ostree: ostreeConfig{
			Name:       "",
			RemoteName: "",
			Ref:        "",
			URL:        "",
		},
		useOstreeRemotes: false,
		defaultSize:      0,
		exports:          []string{},
		requiredPartitionSizes: map[string]datasizes.Size{
			"": 0,
		},
		installWeakDeps:            nil,
		diskImageVPCForceSize:      nil,
		diskImageGiBAligned:        false,
		supportedPartitioningModes: []partition.PartitioningMode{},
		blueprint: blueprintOptions{
			SupportedOptions: []string{},
			RequiredOptions:  []string{},
		},
		arch: &architecture{
			distro: nil,
			arch:   0,
			imageTypes: map[string]distro.ImageType{
				"": nil,
			},
			imageTypeAliases: map[string]string{
				"": "",
			},
		},
		platform: nil,
		image: func(*imageType, *blueprint.Blueprint, distro.ImageOptions, map[string]rpmmd.PackageSet, []rpmmd.RepoConfig, []container.SourceSpec, *rand.Rand) (image.ImageKind, error) {
			panic("not implemented")
		},
		ostreeRef: "",
	}
	return &it, nil
}
