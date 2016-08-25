package rplib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
		Packages           []string
	}
	Udf struct {
		Binary  string
		Command string
	}
	Debug struct {
		Devmode bool
		Ssh     bool
		Xz      bool
	}
	Recovery struct {
		Type                  string // one of "field_transition", "factory_install"
		FsLabel               string `yaml:"filesystem-label"`
		TransitionFsLabel     string
		BootPart              string `yaml:"boot-partition"`
		SystembootPart        string `yaml:"systemboot-partition"`
		WritablePart          string `yaml:"writable-partition"`
		BootImage             string `yaml:"boot-image"`
		SystembootImage       string `yaml:"systemboot-image"`
		WritableImage         string `yaml:"writable-image"`
		SignSerial            bool   `yaml:"sign-serial"`
		SignApiKey            string `yaml:"sign-api-key"`
		SkipFactoryDiagResult string `yaml:"skip-factory-diag-result"`
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

	if config.Yaml.Udf.Command == "" {
		fmt.Println("Error: parse config.yaml failed, need to specify 'udf -> command' field")
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

	if config.Yaml.Debug.Ssh {
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

	io, err := yaml.Marshal(config.Yaml)
	if err != nil {
		panic(err)
	}
	log.Println(string(io))

	// Check if there is any config missing
	errBool := checkConfigs()
	return config, errBool
}

func (configs *ConfigRecovery) ExecuteUDF() {
	args := []string{configs.Yaml.Udf.Command, configs.Yaml.Configs.Release,
		configs.Opt.Store, configs.Yaml.Configs.Store,
		configs.Opt.Device, configs.Yaml.Configs.Device,
		configs.Opt.Channel, configs.Yaml.Configs.Channel,
		configs.Opt.BaseImage, configs.Yaml.Configs.BaseImage,
		configs.Opt.Ssh,
		configs.Opt.Size, configs.Yaml.Configs.Size,
		configs.Opt.Devmode,
		configs.Opt.Kernel, configs.Yaml.Snaps.Kernel,
		configs.Opt.Os, configs.Yaml.Snaps.Os,
		configs.Opt.Gadget, configs.Yaml.Snaps.Gadget}
	for _, snap := range configs.Yaml.Configs.Packages {
		args = append(args, "--install="+snap)
	}
	Shellexec(configs.Yaml.Udf.Binary, args...)
}
