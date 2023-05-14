package core

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var Int = func(s interface{}) int {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return i
}

var Int64 = func(s interface{}) int64 {
	i, _ := strconv.Atoi(fmt.Sprint(s))
	return int64(i)
}

func Itob(i uint64) []byte {
	return []byte(fmt.Sprint(i))
}

func GetUUID() string {
	u, _ := uuid.NewUUID()
	return u.String()
}
func Float64(str interface{}) float64 {
	f, _ := strconv.ParseFloat(fmt.Sprint(str), 64)
	return f
}
func RegistIm(i interface{}) Bucket {
	return NewBucket(regexp.MustCompile("[^/]+$").FindString(reflect.TypeOf(i).PkgPath()))
}

func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 && c != 10 {
			continue
		}
		if c == 127 {
			continue
		}

		dstRunes = append(dstRunes, c)
	}
	return strings.ReplaceAll(string(dstRunes), "￼", "")
}
func GetPush(url string) {

	/*req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("通知出错")
	}
	defer resp.Body.Close()*/
	client := resty.New()
	_, err := client.R().EnableTrace().Get(url)
	if err != nil {
		log.Println("通知出错")
	}
}
