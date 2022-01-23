package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Reader interface {
	Read() (Conf, error)
}

func FromYaml(file string) (Conf, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return Conf{}, err
	}

	var conf Conf
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Conf{}, err
	}

	return conf, nil
}
