package sdk

import (
	"github.com/skratchdot/open-golang/open"
)

// implements UUIDProcessor
type defaultUUIDProcessor struct {
	path string
}

func (dp *defaultUUIDProcessor) ProcessUUID(uuid string) error {
	// 2.``
	path, err := fetchORCodeImage(uuid)

	if err != nil {
		return err
	}
	log.Debugf(`二维码图片地址: %s`, path)

	// 3.
	go func() {
		dp.path = path
		open.Start(path)
	}()
	log.Debug(`请通过微信手机应用程序扫描二维码...`)

	return nil
}

func (dp *defaultUUIDProcessor) UUIDDidConfirm(err error) {
	if len(dp.path) > 0 {
		deleteFile(dp.path)
	}
}
