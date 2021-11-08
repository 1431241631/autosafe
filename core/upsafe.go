package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 响应token字段
type TokenJson struct {
	IdToken          string `json:"idToken"`
	UserNonActivated bool   `json:"userNonActivated"`
	UserNonCompleted bool   `json:"userNonCompleted"`
}

// 服务器响应字段
type ResponseJson struct {
	Code int       `json:"code"`
	Data TokenJson `json:"data"`
}

const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机字符串
func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}

// 登录方法
func Login(username string, password string) (string, error) {
	Url, err := url.Parse("https://token.haust.edu.cn/password/passwordLogin")
	if err != nil {
		fmt.Printf("url parse failed ,err: %s", err.Error())
		return err.Error(), err
	}
	params := url.Values{}
	params.Set("username", username)
	params.Set("password", password)
	params.Set("appId", "com.lantu.MobileCampus.haust")
	params.Set("geo", "")
	params.Set("deviceId", RandChar(24))
	params.Set("osType", "android")

	Url.RawQuery = params.Encode()

	reqUrl := Url.String()

	fmt.Println(reqUrl)

	client := &http.Client{}

	reqPost, _ := http.NewRequest("POST", reqUrl, nil)
	reqPost.Header.Add("Host", "token.haust.edu.cn")
	reqPost.Header.Add("Connection", "Keep-Alive")
	reqPost.Header.Add("User-Agent", "okhttp/3.12.1")

	resp, err := client.Do(reqPost)
	if err != nil {
		fmt.Printf("login err, err: %s", err.Error())
		return err.Error(), err
	}
	defer resp.Body.Close() // 最终关闭流

	bodyContent, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err.Error(), err
	}

	fmt.Println(string(bodyContent))

	var tokenJson ResponseJson
	err = json.Unmarshal(bodyContent, &tokenJson)
	if err != nil {
		fmt.Println("Json err")
		return string(bodyContent), err
	}

	idToken := tokenJson.Data.IdToken

	return idToken, nil

}

func SaveHealthReport(idToken string, healthReport HealthReport) (string, error) {

	jsonData := `{"address":"%s","age":"%d","bodyTemperature":%f,"communityName":null,"communityPhone":null,"createTime":"%s","firstVaccineDate":null,"lastVaccineDate":null,"vaccineType":null,"notVaccineType":null,"notVaccineOtherResult":null,"latitude":"%f","longitude":"%f","homieAddress":null,"isCommunityRemark":null,"isRiskSite":null,"isStayLocal":%d,"isTouchIll":null,"isolationType":null,"journeyAddress":null,"journeyDate":null,"journeyType":null,"onceIllDate":null,"onceTreatHospital":null,"otherSymptomRemark":null,"phone":"%s","recoveryDate":null,"remark":null,"touchDate":null,"touchIllAddress":null,"touchIllDetail":null,"treatDetail":null,"unusualIllDate":null,"unusualSymptomList":[],"vehicleDetail":null,"vehicleType":null,"villageAddress":null,"needUpdate":1,"isAgree":true}`
	timestr := time.Now().Format("2006-01-02 15:04:05")
	jsonData = fmt.Sprintf(jsonData, healthReport.Address, healthReport.Age, healthReport.BodyTemperature, timestr, healthReport.Latitude, healthReport.Longitude, healthReport.IsStayLocal, healthReport.Phone)

	fmt.Println(jsonData)

	client := &http.Client{}

	reqPost, _ := http.NewRequest("POST", "https://yqfkfw.haust.edu.cn/smart-boot/api/healthReport/saveHealthReport", strings.NewReader(jsonData))

	reqPost.Header.Add("Host", "yqfkfw.haust.edu.cn")
	reqPost.Header.Add("Connection", "keep-alive")
	reqPost.Header.Add("X-Terminal-Info", "APP")
	reqPost.Header.Add("X-Device-Info", "APP")
	reqPost.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 9; vivo X21UD A Build/PKQ1.180819.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Mobile Safari/537.36 SuperApp")
	reqPost.Header.Add("content-type", "application/json;charset=UTF-8")
	reqPost.Header.Add("Accept", " */*")
	reqPost.Header.Add("Origin", "https://yqfkfw.haust.edu.cn")
	reqPost.Header.Add("X-Requested-With", "com.lantu.MobileCampus.haust")
	reqPost.Header.Add("Sec-Fetch-Site", "same-origin")
	reqPost.Header.Add("Sec-Fetch-Mode", "cors")
	reqPost.Header.Add("Sec-Fetch-Dest", "empty")
	reqPost.Header.Add("Referer", "https://yqfkfw.haust.edu.cn/serv-h5/")

	reqPost.Header.Add("X-Id-Token", idToken)

	resp, err := client.Do(reqPost)
	if err != nil {
		fmt.Printf("login err, err: %s", err.Error())
		return err.Error(), err
	}
	defer resp.Body.Close() // 最终关闭流

	bodyContent, err := ioutil.ReadAll(resp.Body)

	return string(bodyContent), err
}
