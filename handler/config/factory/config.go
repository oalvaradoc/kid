package factory

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	v2 "git.multiverse.io/eventkit/kit/handler/config/v2"
	"git.multiverse.io/eventkit/kit/log"
	"os"
	"path/filepath"
)

// LoaderFactoryImpl is an implement of config loader
type LoaderFactoryImpl struct{}

// CreateConfigLoader creates a config.Loader for load app configuration
func (l *LoaderFactoryImpl) CreateConfigLoader(version string) (config.Loader, error) {
	var loader config.Loader
	switch version {
	//case "v1":
	//	{
	//		loader = &v1.Loader{}
	//	}
	case "v2":
		{
			loader = &v2.Loader{}
		}
	default:
		{
			return nil, errors.Errorf(constant.SystemInternalError, "cannot found config loader with version:[%s]", version)
		}
	}

	return loader, nil
}

// InitHandlerConfig is used to initialize the config through the specify file path.
func InitHandlerConfig(filePath string, options ...ConfigOption) (configs *config.ServiceConfigs, err error) {
	log.Infosf("start init service config, file path[%s]...", filePath)
	opts := NewConfigOptions()
	for _, o := range options {
		o(&opts)
	}

	factory := &LoaderFactoryImpl{}
	loader, err := factory.CreateConfigLoader(opts.Version)
	if nil != err {
		return nil, err
	}

	cfg, err := loader.LoadConfig(filePath)
	if nil != err {
		return nil, err
	}

	//if true == cfg.Deployment.EnableSecure && "C" == cfg.Deployment.Mode {
	//	filePathElem := strings.Split(filePath, "/")
	//	if len(filePathElem) == 2 {
	//		// delete config file
	//		if filepath, err := GetConfigPath(filePathElem[0], filePathElem[1]); nil == err {
	//			log.Infos("log file path:", filepath)
	//			if err := os.Remove(filepath); nil != err {
	//				panic(err)
	//			}
	//		} else {
	//			panic(err)
	//		}
	//	}
	//}

	if err := InitLocaleLang("conf/locale_*.ini"); nil != err {
		return nil, err
	}

	log.Infosf("successfully to init service config!")
	return cfg, nil
}

// FileExists checks whether the file exists in the disk.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// GetConfigPath joins the file path and file name into absolute file path.
func GetConfigPath(filePath string, fileName string) (string, error) {
	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	workPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	appConfigPath := filepath.Join(workPath, filePath, fileName)
	if !FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, filePath, fileName)
		if !FileExists(appConfigPath) {
			return "", errors.Errorf(constant.SystemInternalError, "cannot found %s/%s", filePath, fileName)
		}
	}

	return appConfigPath, nil
}
