package confmanager

import (
	"os"
	"testing"
)

type MockConfigGetter struct {
	ConfigPath string
	Key        string `json:"key"`
}

func (m *MockConfigGetter) GetConfigPath() string {
	return m.ConfigPath
}

func TestReadConfigFromFile_Success(t *testing.T) {
	// Создаем временный файл конфигурации
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Записываем валидный JSON в файл
	configContent := `{"key": "value"}`
	_, err = tempFile.Write([]byte(configContent))
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Создаем структуру для хранения данных конфигурации
	var cfg map[string]string

	err = readConfigFromFile(tempFile.Name(), &cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем результат
	if cfg["key"] != "value" {
		t.Errorf("expected key 'value', got %s", cfg["key"])
	}
}

func TestReadConfigFromFile_FileNotFound(t *testing.T) {
	// Пытаемся прочитать несуществующий файл
	err := readConfigFromFile("non-existent-file.json", &map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadConfigFromFile_InvalidJSON(t *testing.T) {
	// Создаем временный файл с некорректным JSON
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(`invalid-json`))
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Создаем структуру для хранения данных конфигурации
	var cfg map[string]string

	err = readConfigFromFile(tempFile.Name(), &cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSetCfg_Success(t *testing.T) {
	// Создаем временный файл конфигурации
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	configContent := `{"key": "value"}`
	_, err = tempFile.Write([]byte(configContent))
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Создаем заглушку для ConfigGetter
	mockCfg := &MockConfigGetter{
		ConfigPath: tempFile.Name(),
	}

	// Заглушка для parseCmdCfgFilePath
	parseCmdCfgFilePath := func() {}

	err = SetCfg(mockCfg, parseCmdCfgFilePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mockCfg.Key != "value" {
		t.Errorf("expected key 'value', got %s", mockCfg.Key)
	}
}

func TestSetCfg_FileNotFound(t *testing.T) {
	// Создаем заглушку для ConfigGetter
	mockCfg := &MockConfigGetter{
		ConfigPath: "non-existent-file.json",
	}

	// Заглушка для parseCmdCfgFilePath
	parseCmdCfgFilePath := func() {}

	err := SetCfg(mockCfg, parseCmdCfgFilePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
