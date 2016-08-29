package rplib

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ConfigRecovery struct {
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
		"--size", config.Configs.Size,
		"--kernel", config.Snaps.Kernel,
		"--os", config.Snaps.Os,
		"--gadget", config.Snaps.Gadget}
	if config.Debug.Devmode {
		args = append(args, "--developer-mode")
	}

	if config.Debug.Ssh {
		args = append(args, "--enable-ssh")
	}

	if config.Configs.Store != "" {
		args = append(args, "--store", config.Configs.Store)
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
