package fsloader

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"gopkg.in/yaml.v3"

	"geniusrabbit.dev/adcorelib/storage/loader"
)

var errTargetObjectMustImplementMerge = errors.New(`target object must support "Merge(interface)" method`)

type merger interface {
	Merge(interface{})
}

// PatternLoader returns new FileSystem for some root directory and pattern
func PatternLoader(root, pattern string) loader.LoaderFnk {
	return func(objectTarget interface{}, _ *time.Time) error {
		target, ok := objectTarget.(merger)
		if !ok {
			return errTargetObjectMustImplementMerge
		}

		fileNames, err := filepath.Glob(root + "/" + pattern)
		dataType := reflect.TypeOf(objectTarget).Elem()
		for _, fileName := range fileNames {
			newData := reflect.New(dataType)
			if err := loadFile(fileName, newData.Interface()); err != nil {
				return err
			}
			target.Merge(newData.Interface())
		}
		return err
	}
}

func loadFile(filename string, target interface{}) error {
	ext := filepath.Ext(filename)
	switch ext {
	case ".yml", ".yaml", ".json":
	default:
		return nil // Skip unsupported
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	switch ext {
	case ".yml", ".yaml":
		var interData interface{}
		if err := yaml.Unmarshal(data, &interData); err != nil {
			return err
		}
		if data, err = json.Marshal(interData); err != nil {
			return err
		}
	case ".json":
	default:
		return nil
	}
	return json.Unmarshal(data, target)
}
