package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type UserInfo struct {
	UserName string `yaml:"UserName"`
	PassWord string `yaml:"PassWord"`
	Token    string `yaml:"Token"`
}

type MailConfig struct {
	IsNotify bool   `yaml:"IsNotify"`
	UserName string `yaml:"UserName"`
	AuthCode string `yaml:"AuthCode"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
}

type HealthReport struct {
	Address         string  `yaml:"Address"`
	Latitude        float64 `yaml:"Latitude"`
	Longitude       float64 `yaml:"Longitude"`
	IsStayLocal     int     `yaml:"IsStayLocal"`
	Age             int     `yaml:"Age"`
	BodyTemperature float32 `yaml:"BodyTemperature"`
	Phone           string  `yaml:"Phone"`
}

type Config struct {
	UserInfo     UserInfo     `yaml:"UserInfo"`
	Mail         MailConfig   `yaml:"Mail"`
	HealthReport HealthReport `yaml:"HealthReport"`
}

var ConfigDate Config

// 封装下发送邮件，这里调用次数多
func SendMailPlus(subject string, body string) {
	if ConfigDate.Mail.IsNotify {
		sendMail(ConfigDate.Mail.UserName, subject, body, ConfigDate.Mail.AuthCode, ConfigDate.Mail.Host, ConfigDate.Mail.Port)
	}
}

// 读取配置文件(有条件的可以改成数据库)
func ReadConfig() error {
	file, err := os.Open("./config.yaml")
	if err != nil {
		fmt.Printf("file open error %s", err.Error())
		return err
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	err = yaml.Unmarshal(data, &ConfigDate)
	if err != nil {
		return err
	}

	// 这里持久化是为了防止每次都登录，实际上token有效期有一个月。实测云函数计算的持久化并不是很有效，不太明白调用规则(目测每次调用并不是在同一实例中)。好在能正常运行
	data, err = ioutil.ReadFile("/tmp/token")
	if err != nil {
		ConfigDate.UserInfo.Token = ""
	} else {
		ConfigDate.UserInfo.Token = string(data)
	}

	return nil
}

func HandleRequest_() (string, error) {

	err := ReadConfig()
	if err != nil {
		return err.Error(), err
	}
	if len(ConfigDate.UserInfo.Token) < 10 {
		token, err := Login(ConfigDate.UserInfo.UserName, ConfigDate.UserInfo.PassWord)
		fmt.Println(token)
		if err != nil {
			SendMailPlus("登录失败", token)
			return "登录失败", err
		}

		ConfigDate.UserInfo.Token = token
		// 保存一下token
		ioutil.WriteFile("/tmp/token", []byte(ConfigDate.UserInfo.Token), 0777)
		SendMailPlus("获取Token完成", token)
	}

	res, err := SaveHealthReport(ConfigDate.UserInfo.Token, ConfigDate.HealthReport)
	if err != nil {
		SendMailPlus("已为您自动打卡,但是似乎有什么问题", res+err.Error())
		return "打卡完成", err
	}
	SendMailPlus("已为您自动打卡", res)

	return "打卡完成", nil
}
