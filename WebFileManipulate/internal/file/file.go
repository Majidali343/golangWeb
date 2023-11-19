// internal/file/file.go
package file

import "io/ioutil"


func ReadFile(filePath string) ([]byte, error) {
    return ioutil.ReadFile(filePath)
}
