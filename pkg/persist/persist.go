package persist

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var lock sync.Mutex

func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return errors.WithStack(err)
	}
	r := bytes.NewReader(b)

	_, err = io.Copy(f, r)
  	return errors.WithStack(err)
  }

func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Open(path)
	if err != nil {
	  return errors.WithStack(err)
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
  }