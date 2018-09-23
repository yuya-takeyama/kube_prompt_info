package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

type kubeContext struct {
	Namespace string `yaml:"namespace"`
}

type kubeContextWithName struct {
	Context kubeContext `yaml:"context"`
	Name    string      `yaml:"name"`
}

type kubeConfig struct {
	APIVersion     string                `yaml:"apiVersion"`
	CurrentContext string                `yaml:"current-context"`
	Contexts       []kubeContextWithName `yaml:"contexts"`
}

func main() {
	home, homedirErr := homedir.Dir()
	if homedirErr != nil {
		fatalf("failed to get homedir")
	}

	configPath := filepath.Join(home, ".kube", "config")
	file, openErr := os.Open(configPath)
	if openErr != nil {
		fatalf("failed to open kubectl config: %s", configPath)
	}

	buf, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		fatalf("failed to read kubectl config: %s", configPath)
	}

	var config kubeConfig
	yamlErr := yaml.Unmarshal(buf, &config)
	if yamlErr != nil {
		fatalf("failed to parse kubectl config: %s", configPath)
	}

	var namespace string
	for _, ctxWithName := range config.Contexts {
		if ctxWithName.Name == config.CurrentContext {
			namespace = ctxWithName.Context.Namespace
		}
	}

	fmt.Printf("%s%s/%s%s", os.Getenv("KUBE_PROMPT_INFO_PREFIX"), config.CurrentContext, namespace, os.Getenv("KUBE_PROMPT_INFO_SUFFIX"))
}

func fatalf(message string, params ...interface{}) {
	if os.Getenv("KUBE_PROMPT_INFO_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "kube-prompt-info: "+message+"\n", params...)
	}
	os.Exit(1)
}
