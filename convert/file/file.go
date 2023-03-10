package file

import (
	"os"
)

type Data struct {
	data []byte
}

func New(data []byte) *Data {
	return &Data{data}
}

func (f *Data) Save(filename string) error {
	err := os.WriteFile(filename, f.data, 0644)
	return err
}
