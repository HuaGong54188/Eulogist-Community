package Eulogist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

// FileExist 检查 path 对应路径的文件是否存在。
// 如果不存在，或该路径指向一个文件夹，则返回假，否则返回真
func FileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

// WriteJsonFile 将 content 以 JSON
// 形式写入到 path 指代的文件处
func WriteJsonFile(path string, content any) error {
	contentBytes, _ := json.Marshal(content)
	// 将内容转换为 JSON 格式
	buffer := bytes.NewBuffer([]byte{})
	json.Indent(buffer, contentBytes, "", "	")
	// 格式化 JSON
	err := os.WriteFile(path, buffer.Bytes(), 0600)
	if err != nil {
		return fmt.Errorf("WriteJsonFile: %v", err)
	}
	// 将 JSON 写入文件
	return nil
}

// ReadEulogistConfig 在当前目录读取赞颂者的配置文件。
// 如果没有对应的文件，则将尝试生成默认配置文件。
//
// 生成默认配置文件期间需要从控制台读取用户输入，
// 读取的内容包括需要进入的租赁服号及其密码，
// 以及 FastBuilder 原生验证服务器的 Token。
func ReadEulogistConfig() (*EulogistConfig, error) {
	var cfg EulogistConfig

	if !FileExist("eulogist_config.json") {
		config, err := GenerateEulogistConfig()
		if err != nil {
			return nil, fmt.Errorf("ReadEulogistConfig: %v", err)
		}
		return config, nil
	}

	fileBytes, err := os.ReadFile("eulogist_config.json")
	if err != nil {
		return nil, fmt.Errorf("ReadEulogistConfig: %v", err)
	}

	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("ReadEulogistConfig: %v", err)
	}

	return &cfg, nil
}

// GenerateEulogistConfig 在当前目录生成赞颂者的默认配置文件，
// 并返回该配置文件。此函数会从控制台读取用户输入，
// 读取的内容包括需要进入的租赁服号及其密码，
// 以及 FastBuilder 原生验证服务器的 Token
func GenerateEulogistConfig() (config *EulogistConfig, err error) {
	cfg := EulogistConfig{}

	pterm.Info.Printf("输入租赁服号: ")
	fmt.Scanln(&cfg.RentalServerCode)

	pterm.Info.Printf("输入租赁服密码(没有可留空): ")
	fmt.Scanln(&cfg.RentalServerPassword)

	pterm.Info.Printf("输入FB原生验证或者咕咕验证的token: ")
	fmt.Scanln(&cfg.FBToken)

	err = WriteJsonFile("eulogist_config.json", cfg)
	if err != nil {
		return nil, fmt.Errorf("GenerateEulogistConfig: %v", err)
	}

	return &cfg, nil
}