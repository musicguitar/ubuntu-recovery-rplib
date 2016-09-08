package rplib

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ConfigRecovery struct {
	// TODO: deprecate Snaps parameter
	Project string
	Snaps   struct {
		Kernel string
		Os     string
		Gadget string
	}
	Configs struct {
		// TODO: deprecate Store parameter
		Arch               string
		BaseImage          string
		RecoveryType       string
		RecoverySize       string
		Release            string
		Store              string
		Device             string // parameter for ubuntu-device-flash
		Channel            string
		Size               string
		OemPreinstHookDir  string `yaml:"oem-preinst-hook-dir"`
		OemPostinstHookDir string `yaml:"oem-postinst-hook-dir"`
		OemLogDir          string
		Packages           []string
		PartitionType      string `yaml:"partition-type"`
		Bootloader         string `yaml:"bootloader"`
		ModelAssertion     string
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

func (config *ConfigRecovery) checkConfigs() (err error) {
	log.Println("check configs ... ")

	if config.Project == "" {
		err = errors.New("'project' field not presented")
		log.Println(err)
	}

	if config.Snaps.Kernel == "" {
		err = errors.New("'snaps -> kernel' field not presented")
		log.Println(err)
	}

	if config.Snaps.Os == "" {
		err = errors.New("'snaps -> os' field not presented")
		log.Println(err)
	}

	if config.Snaps.Gadget == "" {
		err = errors.New("'snaps -> gadget' field not presented")
		log.Println(err)
	}

	if config.Configs.Arch == "" {
		err = errors.New("'configs -> arch' field not presented")
		log.Println(err)
	} else if config.Configs.Arch != "amd64" && config.Configs.Arch != "arm" && config.Configs.Arch != "arm64" {
		err = errors.New("'recovery -> Arch' only accept \"amd64\" or \"arm\" or \"arm64\"")
		log.Println(err)
	}

	if config.Configs.BaseImage == "" {
		err = errors.New("'configs -> baseimage' field not presented")
		log.Println(err)
	}

	if config.Configs.RecoveryType == "" {
		err = errors.New("'configs -> recoverytype' field not presented")
		log.Println(err)
	}

	if config.Configs.RecoverySize == "" {
		err = errors.New("'configs -> recoverysize' field not presented")
		log.Println(err)
	}

	if config.Configs.Release == "" {
		err = errors.New("'configs -> release' field not presented")
		log.Println(err)
	}

	if config.Configs.Channel == "" {
		err = errors.New("'configs -> channel' field not presented")
		log.Println(err)
	}

	if config.Configs.Size == "" {
		err = errors.New("'configs -> size' field not presented")
		log.Println(err)
	}

	if config.Configs.PartitionType == "" {
		err = errors.New("'recovery -> PartitionType' field not presented")
		log.Println(err)
	} else if config.Configs.PartitionType != "gpt" && config.Configs.PartitionType != "mbr" {
		err = errors.New("'recovery -> PartitionType' only accept \"gpt\" or \"mbr\"")
		log.Println(err)
	}

	if config.Configs.Bootloader == "" {
		err = errors.New("'recovery -> PartitionType' field not presented")
		log.Println(err)
	} else if config.Configs.Bootloader != "grub" && config.Configs.Bootloader != "u-boot" {
		err = errors.New("'recovery -> PartitionType' only accept \"grub\" or \"u-boot\"")
		log.Println(err)
	}

	if config.Udf.Binary == "" {
		err = errors.New("'udf -> binary' field not presented")
		log.Println(err)
	}

	if config.Udf.Command == "" {
		err = errors.New("'udf -> command' field not presented")
		log.Println(err)
	}

	if config.Recovery.FsLabel == "" {
		err = errors.New("'recovery -> filesystem-label' field not presented")
		log.Println(err)
	}

	return err
}

func (config *ConfigRecovery) Load(configFile string) error {
	log.Println("Loading config file %s ...", configFile)
	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		return err
	}

	// Parse config file and store in configs
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	// Check if there is any config missing
	err = config.checkConfigs()
	return err
}

func (config *ConfigRecovery) ExecuteUDF() {
	args := []string{
		config.Udf.Command, config.Configs.Release,
		"--channel", config.Configs.Channel,
		"--output", config.Configs.BaseImage,
		config.Configs.ModelAssertion,
	}
	if config.Debug.Devmode {
		args = append(args, "--developer-mode")
	}

	if config.Debug.Ssh {
		args = append(args, "--enable-ssh")
	}

	if config.Configs.Device != "" {
		args = append(args, "--device", config.Configs.Device)
	}

	for _, snap := range config.Configs.Packages {
		args = append(args, "--install="+snap)
	}
	Shellexec(config.Udf.Binary, args...)
}

func (config *ConfigRecovery) String() string {
	io, err := yaml.Marshal(*config)
	if err != nil {
		panic(err)
	}
	return string(io)
}
