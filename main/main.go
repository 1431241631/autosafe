package main

import (
	"autosafe/core"

	"github.com/aliyun/fc-runtime-go-sdk/fc"
)

func main() {
	fc.Start(HandleRequest)
}

func HandleRequest(event map[string]interface{}) (string, error) {
	return core.HandleRequest_()
}
