package entity

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
)

func NewRequest(root, collection, name string) error {
	collectionDir := root
	if collection != Global {
		collectionDir = path.Join(root, collection)
		_, err := os.Stat(collectionDir)
		if err != nil {
			return errors.Errorf("collection %s does not exist", collection)
		}
	}

	requestPath := path.Join(collectionDir, name+".json")
	_, err := os.Stat(requestPath)
	if err == nil {
		println("request", name, "already exists")
		return nil
	}

	requestConfig, err := json.Marshal(&RequestConfig{})
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}

	err = os.WriteFile(requestPath, requestConfig, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}

	return nil
}

func RemoveRequest(root, collection, name string) error {
	requestPath := path.Join(root, name+".json")
	if collection != Global {
		requestPath = path.Join(root, collection, name+".json")
	}
	err := os.Remove(requestPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return errors.Wrap(err, "could not delete request")
}

func EditRequest(root, collection, name string, opts ...ReqOpt) error {
	workingDir := root
	if collection != Global {
		workingDir = path.Join(root, collection)
	}
	requestPath := path.Join(workingDir, name+".json")
	_, err := os.Stat(requestPath)
	if err != nil {
		println("request", name, "not found")
		return nil
	}

	cfg := viper.New()
	cfg.SetConfigType("json")
	cfg.SetConfigFile(requestPath)

	err = cfg.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "could not read request config")
	}

	reqConfig := RequestConfig{}
	err = cfg.Unmarshal(&reqConfig)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal request config")
	}

	for _, opt := range opts {
		reqConfig = opt(reqConfig)
	}

	requestConfig, err := json.Marshal(&reqConfig)
	if err != nil {
		return errors.Wrap(err, "could not edit request")
	}

	err = os.WriteFile(requestPath, requestConfig, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not edit request")
	}

	return nil
}

func ExecuteRequest(root, collection, name string) error {
	workingDir := root
	cfg := newDefaultConfig()
	cfg.SetConfigType("json")
	cfg.SetConfigFile(path.Join(workingDir, CfgFileNameJSON))
	err := cfg.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "could not read global config")
	}

	if collection != Global {
		workingDir = path.Join(root, collection)
		cfg.SetConfigFile(path.Join(workingDir, CfgFileNameJSON))
		err := cfg.ReadInConfig()
		if err != nil {
			return errors.Wrap(err, "could not read collection config")
		}
	}

	cfg.SetConfigFile(path.Join(workingDir, name+".json"))
	err = cfg.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "could not read request config")
	}

	reqConfig := RequestConfig{}
	err = cfg.Unmarshal(&reqConfig)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal request config")
	}

	cmd := exec.Command("curl", reqConfig.toCurl()...)
	fmt.Printf("Executing:\n%s\n", cmd)
	out, err := cmd.CombinedOutput()
	println(string(out))
	return err
}
