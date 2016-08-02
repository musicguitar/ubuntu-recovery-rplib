package rplib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type OptString struct {
	Kernel    string
	Os        string
	Gadget    string
	BaseImage string
	Store     string
	Device    string
	Channel   string
	Size      string
	Devmode   string
	Ssh       string
}

// config.yaml from yaml file
type YamlConfig struct {
	Project string
	Snaps   struct {
		Kernel string
		Os     string
		Gadget string
	}
	Configs struct {
		Arch               string
		BaseImage          string
		RecoveryType       string
		RecoverySize       string
		Release            string
		Store              string
		Device             string
		Channel            string
		Size               string
		OemPreinstHookDir  string `yaml:"oem-preinst-hook-dir"`
		OemPostinstHookDir string `yaml:"oem-postinst-hook-dir"`
		OemLogDir          string
	}
	Udf struct {
		Binary string
		Option string
	}
	Debug struct {
		Devmode bool
		Ssh     bool
		Xz      bool
	}
	Recovery struct {
		FsLabel         string `yaml:"filesystem-label"`
		BootPart        string `yaml:"boot-partition"`
		SystembootPart  string `yaml:"systemboot-partition"`
		WritablePart    string `yaml:"writable-partition"`
		BootImage       string `yaml:"boot-image"`
		SystembootImage string `yaml:"systemboot-image"`
		WritableImage   string `yaml:"writable-image"`
		SignSerial      bool   `yaml:"sign-serial"`
		SignApiKey      string `yaml:"sign-api-key"`
	}
}

type ConfigRecovery struct {
	Yaml YamlConfig
	Opt  OptString
}

var config ConfigRecovery

func loadDefaultOptValue() {
	config.Opt.Kernel = "--kernel"
	config.Opt.Os = "--os"
	config.Opt.Gadget = "--gadget"
	config.Opt.BaseImage = "--output"
	config.Opt.Store = ""
	config.Opt.Device = ""
	config.Opt.Channel = "--channel"
	config.Opt.Size = "--size"
	config.Opt.Devmode = ""
	config.Opt.Ssh = ""
}

func checkConfigs() bool {
	fmt.Printf("check configs ... \n")

	errCount := 0
	if config.Yaml.Project == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'project' field")
		errCount++
	}

	if config.Yaml.Snaps.Kernel == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'snaps -> kernel' field")
		errCount++
	}

	if config.Yaml.Snaps.Os == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'snaps -> os' field")
		errCount++
	}

	if config.Yaml.Snaps.Gadget == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'snaps -> gadget' field")
		errCount++
	}

	if config.Yaml.Configs.BaseImage == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> baseimage' field")
		errCount++
	}

	if config.Yaml.Configs.RecoveryType == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> recoverytype' field")
		errCount++
	}

	if config.Yaml.Configs.RecoverySize == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> recoverysize' field")
		errCount++
	}

	if config.Yaml.Configs.Release == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> release' field")
		errCount++
	}

	if config.Yaml.Configs.Channel == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> channel' field")
		errCount++
	}

	if config.Yaml.Configs.Size == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'config.yaml.-> size' field")
		errCount++
	}

	if config.Yaml.Udf.Binary == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'udf -> binary' field")
		errCount++
	}

	if config.Yaml.Udf.Option == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'udf -> option' field")
		errCount++
	}

	if config.Yaml.Recovery.FsLabel == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'recovery -> filesystem-label' field")
		errCount++
	}

	if errCount > 0 {
		return true
	}

	if config.Yaml.Debug.Devmode {
		config.Opt.Devmode = "--developer-mode"
	}

	if config.Yaml.Debug.Devmode {
		config.Opt.Ssh = "--enable-ssh"
	}

	if config.Yaml.Configs.Store != "" {
		config.Opt.Store = "--store"
	}

	if config.Yaml.Configs.Device != "" {
		config.Opt.Device = "--device"
	}

	return false
}

func printConfigs() {
	fmt.Printf("Configs from yaml file\n")
	fmt.Println("-----------------------------------------------")
	fmt.Println("project: ", config.Yaml.Project)
	fmt.Println("kernel: ", config.Yaml.Snaps.Kernel)
	fmt.Println("os: ", config.Yaml.Snaps.Os)
	fmt.Println("gadget: ", config.Yaml.Snaps.Gadget)
	fmt.Println("baseimage: ", config.Yaml.Configs.BaseImage)
	fmt.Println("recoverytype: ", config.Yaml.Configs.RecoveryType)
	fmt.Println("recoverysize: ", config.Yaml.Configs.RecoverySize)
	fmt.Println("release: ", config.Yaml.Configs.Release)
	fmt.Println("store: ", config.Yaml.Configs.Store)
	fmt.Println("device: ", config.Yaml.Configs.Device)
	fmt.Println("channel: ", config.Yaml.Configs.Channel)
	fmt.Println("size: ", config.Yaml.Configs.Size)
	fmt.Println("oem-preinst-hook-dir: ", config.Yaml.Configs.OemPreinstHookDir)
	fmt.Println("oem-postinst-hook-dir: ", config.Yaml.Configs.OemPostinstHookDir)
	fmt.Println("oemlogdir: ", config.Yaml.Configs.OemLogDir)
	fmt.Println("udf binary: ", config.Yaml.Udf.Binary)
	fmt.Println("udf option: ", config.Yaml.Udf.Option)
	fmt.Println("devmode: ", config.Yaml.Debug.Devmode)
	fmt.Println("ssh: ", config.Yaml.Debug.Ssh)
	fmt.Println("xz: ", config.Yaml.Debug.Xz)
	fmt.Println("fslabel: ", config.Yaml.Recovery.FsLabel)
	fmt.Println("boot partition: ", config.Yaml.Recovery.BootPart)
	fmt.Println("system-boot partition: ", config.Yaml.Recovery.SystembootPart)
	fmt.Println("writable partition: ", config.Yaml.Recovery.WritablePart)
	fmt.Println("boot image: ", config.Yaml.Recovery.BootImage)
	fmt.Println("system-boot image: ", config.Yaml.Recovery.SystembootImage)
	fmt.Println("writable image: ", config.Yaml.Recovery.WritableImage)
	fmt.Println("sign serial: ", config.Yaml.Recovery.SignSerial)
	fmt.Println("sign api key: ", config.Yaml.Recovery.SignApiKey)
	fmt.Println("-----------------------------------------------")
}

func LoadYamlConfig(configFile string) (ConfigRecovery, bool) {
	fmt.Printf("Loading config file %s ...\n", configFile)
	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)

	// Load default option string
	loadDefaultOptValue()

	if err != nil {
		fmt.Printf("Error: can not load %s\n", configFile)
		panic(err)
	}

	// Parse config.yaml and store in configs
	err = yaml.Unmarshal(yamlFile, &config.Yaml)
	if err != nil {
		fmt.Printf("Error: parse %s failed\n", configFile)
		panic(err)
	}

	// Check if there is any config missing
	errBool := checkConfigs()
	printConfigs()
	return config, errBool
}
