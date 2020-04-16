package audit

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
	apiUrl     = "http://api.botsmart.cn/v1/check/send"
	appKey      = "f7efe97f915c4bb6b023afa69ed03d89"  //产品密钥ID
	appSecret  = ""                                  //产品密钥
	businessId = ""                                 //业务ID
)

func main() {
	params := url.Values{}

	ret := check(params)
	fmt.Println(ret)
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
	s := hex.EncodeToString(h.Sum(nil) )
	return s
}


func check(params url.Values) *simplejson.Json {
	params["app_id"] = []string{appKey}
	params["business_id"] = []string{businessId}
	params["unique_id"] = []string{uuid.NewV1().String()}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["data"] = []string{"测试内容"}
	params["signature"] = []string{genSignature(params)}

	resp, err := http.PostForm(apiUrl,  params)

	if err != nil {
		fmt.Println("调用API接口失败:", err)
		return nil
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
}
