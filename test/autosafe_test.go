package TEST

import (
	"autosafe/core"
	"testing"
)

func Test_login(t *testing.T) {
	core.ReadConfig()
	token, err := core.Login(core.ConfigDate.UserInfo.UserName, core.ConfigDate.UserInfo.PassWord)
	t.Log("res")
	t.Log(token, err)
}

func Test_saveHealthReport(t *testing.T) {
	core.ReadConfig()
	res, _ := core.SaveHealthReport("token", core.ConfigDate.HealthReport)
	t.Log(res)
}

func Test_sendMail(t *testing.T) {
	core.ReadConfig()
	core.SendMailPlus("TEST", "这是一条测试信息")
}

func Test_handler(t *testing.T) {
	core.HandleRequest_()
}
