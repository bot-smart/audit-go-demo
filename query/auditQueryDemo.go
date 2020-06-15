package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	queryUrl   = "https://api.botsmart.cn/v1/check/query"
	appKey     = "" //产品密钥ID
	appSecret  = "" //产品密钥
	businessId = "" //业务ID
)

func main() {
	taskIds := []int64{1272376584721129474, 1266367504416641025}
	jsonString, _ := json.Marshal(taskIds)
	params := url.Values{"taskIds": []string{string(jsonString)}}
	ret := query(params)
	fmt.Println(ret)
}

func query(params url.Values) *simplejson.Json {
	params["app_id"] = []string{appKey}
	params["business_id"] = []string{businessId}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["signature"] = []string{genSignature(params)}

	resp, err := http.PostForm(queryUrl, params)

	if err != nil {
		fmt.Println("调用API接口失败:", err)
		return nil
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
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
