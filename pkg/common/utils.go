/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package common

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

var config Config

// Service YAML Config Structure
type Service struct {
	Enable             bool              `yaml:"enable"`
	PodmanId           string            `yaml:"podman_id"`
	PodmanImage        string            `yaml:"podman_image"`
	PodmanName         string            `yaml:"podman_name"`
	PodName            string            `yaml:"pod_name"`
	ContainerName      string            `yaml:"container_name"`
	StrictPodNameMatch bool              `yaml:"strict_pod_name_match"`
	Path               []string          `yaml:"path"`
	Hosts              []string          `yaml:"hosts"`
	ServiceCommand     string            `yaml:"service_command"`
	CatOutput          bool              `yaml:"cat_output"`
	ConfigMapping      map[string]string `yaml:"config_mapping"`
}

type Config struct {
	Services map[string]Service `yaml:"services"`
}

// Shell execution functions:
func ExecCmd(cmd string) ([]string, error) {
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return strings.Split(string(output), "\n"), err
	}
	return strings.Split(string(output), "\n"), nil
}

func ExecCmdSimple(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func TestOCConnection() bool {
	cmd := "oc whoami"
	_, err := ExecCmd(cmd)
	if err != nil {
		return false
	}
	return true
}

func TestSshConnection(sshCmd string) bool {
	cmd := sshCmd + " ls"
	_, err := ExecCmd(cmd)
	if err != nil {
		return false
	}
	return true
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetNestedFieldValue(data interface{}, keyName string) interface{} {
	val := reflect.ValueOf(data)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	field := val.FieldByName(keyName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

func LoadServiceConfigFile(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding YAML:", err)
		return config, err
	}
	return config, nil
}

func ConvertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case []string:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var result string
	for _, part := range parts {
		result += strings.Title(part)
	}
	return result
}
