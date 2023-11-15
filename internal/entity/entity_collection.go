package entity

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"path"
)

func NewCollection(root, name string) error {
	collectionPath := path.Join(root, name)

	_, err := os.Stat(collectionPath)
	if err == nil {
		println("collection", name, "already exists")
		return nil
	}

	err = os.Mkdir(collectionPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not create collection")
	}

	collectionCfg, err := json.Marshal(&CollectionConfig{})
	if err != nil {
		return errors.Wrap(err, "could not create collection")
	}

	err = os.WriteFile(path.Join(collectionPath, CfgFileNameJSON), collectionCfg, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not create collection")
	}

	return nil
}

func RemoveCollection(root, name string) error {
	collectionPath := path.Join(root, name)
	err := os.RemoveAll(collectionPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return errors.Wrap(err, "could not delete collection")
}
