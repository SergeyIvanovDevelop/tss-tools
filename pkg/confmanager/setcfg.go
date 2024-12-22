package confmanager

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigGetter interface {
	GetConfigPath() string
}

func SetCfg(cfgPtr ConfigGetter, parseCmdCfgFilePath func()) error {
	parseCmdCfgFilePath()
	parseEnvVariables(cfgPtr)

	cfgFilePath := cfgPtr.GetConfigPath() // Либо из ENV либо из default
	if cfgFilePath != "" {
		err := readConfigFromFile(cfgFilePath, cfgPtr)
		if err != nil {
			return fmt.Errorf("error reading config from file '%s': %w", cfgFilePath, err)
		}
	}

	ParseCmdFlagsSecondary()
	parseEnvVariables(cfgPtr)

	return nil
}

func readConfigFromFile(filePath string, cfg any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return nil
}
