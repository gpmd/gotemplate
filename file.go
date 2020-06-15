package gotemplate

import (
	"io/ioutil"
	"os"

	"github.com/spf13/afero"
)

var fs afero.Fs

// RegisterFS (afero) virtual filesystem for ProcessTemplateFile
func RegisterFS(filesystem afero.Fs) {
	fs = filesystem
}

// ProcessTemplateFile processes golang template file
func ProcessTemplateFile(template string, bundle interface{}) ([]byte, error) {
	var byteValue []byte
	var err error
	if fs == nil {
		tf, err := os.Open(template)
		if err != nil {
			return nil, err
		}
		defer tf.Close()
		byteValue, err = ioutil.ReadAll(tf)
	} else {
		tf, err := fs.Open(template)
		if err != nil {
			return nil, err
		}
		byteValue, err = ioutil.ReadAll(tf)
		defer tf.Close()
	}
	if err != nil {
		return nil, err
	}
	output, err := Template(string(byteValue), bundle)
	if err != nil {
		return []byte{}, err
	}
	return []byte(output), nil
}

// MustProcessTemplateFile processes golang template file otherwise panics
func MustProcessTemplateFile(template string, bundle interface{}) string {
	tf, err := os.Open(template)
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(tf)
	output, _ := Template(string(byteValue), bundle)
	defer tf.Close()
	return output
}
