package tools

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

func DecodeGBK(s []byte) ([]byte, error) {
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(bytes.NewReader(s))

	return ioutil.ReadAll(reader)
}

func FileWrite(file string, content []byte) error {

	filepath := path.Join(file)

	basepath := path.Dir(filepath)
	// 检测文件夹是否存在   若不存在  创建文件夹
	if _, err := os.Stat(basepath); err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(basepath, os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		return err
	}

	_, err = f.Write(content)

	return err
}

func FileAppend(file, content string) error {

	filepath := path.Join(file)

	basepath := path.Dir(filepath)
	// 检测文件夹是否存在   若不存在  创建文件夹
	if _, err := os.Stat(basepath); err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(basepath, os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	data := []byte(content)

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)

	if err != nil {
		return err
	}

	_, err = f.Write(data)

	return err
}

func RandomDelay(seed int64) time.Duration {
	return time.Duration(rand.Int63n(seed+1000)) * time.Millisecond
}
