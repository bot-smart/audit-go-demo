package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	apiUrl     = "https://api.botsmart.cn/v1/check/send"
	resultUrl  = "https://api.botsmart.cn/v1/check/query"
	appKey     = "" //产品密钥ID
	appSecret  = "" //产品密钥
	businessId = "" //业务ID
)

func main() {
	// 调用检测接口
	ret := check()
	fmt.Println(ret)

	// 获取taskId
	taskId, _ := ret.Get("data").Get("taskId").String()
	fmt.Println(taskId)

	// 根据taskId取结果
	ret2 := getResult(taskId)
	fmt.Println(ret2)
}

//获取uuid
func getUuid() string {
	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return ""
	}
	return u2.String()
}

//生成签名信息
func genSignature(params url.Values) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + "=" + params[key][0] + "&"
	}
	str := paramStr[0 : len(paramStr)-1]
	str += appSecret
	h := sha1.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		return ""
	}
	s := hex.EncodeToString(h.Sum(nil))
	return s
}

func check() *simplejson.Json {
	params := url.Values{}
	params["app_id"] = []string{appKey}
	params["business_id"] = []string{businessId}
	params["unique_id"] = []string{getUuid()}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["data"] = []string{"测试内容"}
	params["signature"] = []string{genSignature(params)}

	resp, err := http.PostForm(apiUrl, params)

	if err != nil {
		fmt.Println("调用API接口失败:", err)
		return nil
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
}

func getResult(taskId string) *simplejson.Json {
	params := url.Values{}
	params["app_id"] = []string{appKey}
	params["business_id"] = []string{businessId}
	params["taskIds"] = []string{"[" + taskId + "]"}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["signature"] = []string{genSignature(params)}

	resp, err := http.PostForm(resultUrl, params)

	if err != nil {
		fmt.Println("调用API接口失败:", err)
		return nil
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
}
