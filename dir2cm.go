package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type MetaData struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ConfigMap struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   MetaData          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

func EmptyConfigMap(name string, namespace string) *ConfigMap {
	cm := &ConfigMap{
		ApiVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: MetaData{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}
	return cm
}

func ParseTemplateYaml(fpath string) (*ConfigMap, error) {
	contents, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	c := ConfigMap{}
	err = yaml.Unmarshal([]byte(contents), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Adds a file
func (c *ConfigMap) AddFile(f *ConfigMapFile) error {
	c.Data[f.Name] = string(f.Contents)
	return nil
}

func (c *ConfigMap) DumpYaml() {
	yml, err := yaml.Marshal(*c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(yml))
}

/**
 * ConfigMapFile
 *
 * Will be come a key in the ConfigMap
 */
type ConfigMapFile struct {
	// Actual FS path
	Path string
	// Name/key for configmap (basename(Path))
	Name string `yaml:"name"`
	// Contents (as bytes)
	Contents []byte
}

func NewConfigMapFile(fpath string) (*ConfigMapFile, error) {
	contents, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	cm := &ConfigMapFile{
		Path:     fpath,
		Name:     path.Base(fpath),
		Contents: contents,
	}
	return cm, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var name string
	var dir string
	var namespace string
	var template string

	flag.StringVar(&name, "name", "my-config", "The ConfigMap Metadata.Name")
	flag.StringVar(&name, "n", "my-config", "The ConfigMap Metadata.Name")
	flag.StringVar(&dir, "dir", cwd, "The input directory")
	flag.StringVar(&dir, "d", cwd, "The input directory")
	flag.StringVar(&namespace, "namespace", "", "The ConfigMap Metadata.Namespace")
	flag.StringVar(&namespace, "ns", "", "The ConfigMap Metadata.Namespace")
	flag.StringVar(&template, "template", "", "The ConfigMap file to use as a template")
	flag.StringVar(&template, "t", "", "The ConfigMap file to use as a template")

	usage := fmt.Sprintf(`Usage of ./dir2cm:
	-d, --dir string
		  The input directory (default "%s")
	-n, --name string
		  The ConfigMap Metadata.Name (default "my-config")
	-ns, --namespace string
		  The ConfigMap Metadata.Namespace
	-t, --template string
		  The ConfigMap file to use as a template
	`, cwd)
	flag.Usage = func() { fmt.Print(usage) }

	flag.Parse()

	//var files []string
	files, err := os.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	var cm *ConfigMap
	if template != "" {
		cm, err = ParseTemplateYaml(template)
		if err != nil {
			log.Warn(err)
			cm = EmptyConfigMap(name, namespace)
		}
	} else {
		cm = EmptyConfigMap(name, namespace)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullpath := path.Join(dir, file.Name())
		cmf, err := NewConfigMapFile(fullpath)
		if err != nil {
			log.Warnf("Problem with file: %s", err)
			continue
		}
		err = cm.AddFile(cmf)
		if err != nil {
			log.Warnf("Couldn't add file %s (%s)", fullpath, err)
		}
	}

	cm.DumpYaml()
}
