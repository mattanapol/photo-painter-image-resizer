package image_processing

import (
	"log"
	"os"
	"path"
)

func WriteFile(filePath, fileName string, content []byte) error {

	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	out, err := os.Create(path.Join(filePath, fileName))
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Println(err)
		}
	}(out)
	_, err = out.Write(content)
	if err != nil {
		return err
	}
	return nil
}
