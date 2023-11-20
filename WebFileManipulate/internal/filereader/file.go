// internal/file/file.go
package filereader

import "io/ioutil"

func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
