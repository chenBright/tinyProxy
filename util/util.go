package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// 将 host、port 转换成 url 字符串
func HostPortToAddress(host string, port uint16) string {
	return host + ":" + strconv.Itoa(int(port))
}

// 获取 url 中的 host
func UrlToHost(url string) string {
	return strings.Split(url, ":")[0]
}

// 获取绝对路径
func AbsolutePath(relpath string) string {
	absolutePath, err := filepath.Abs(relpath)
	if err != nil {
		log.Println("current path error: %s", err)
	}

	return absolutePath
}

// 获取主目录的绝对路径
func HomePath() string {
	return AbsolutePath(".")
}

// 获取 element 在 slice 中的索引
func SliceIndex(slice interface{}, element interface{}) int {
	var index int = -1

	// 通过反射来获取变量的类型
	sliceValue := reflect.ValueOf(slice)
	// 如果 slice 的类型不是 Slice，则返回 -1
	if sliceValue.Kind() != reflect.Slice {
		return index
	}

	elementValue := reflect.ValueOf(element).Interface()
	length := sliceValue.Len()
	for i := 0; i < length; i++ {
		// 通过反射来获取指定索引的元素，并进行比较
		itemValue := sliceValue.Index(i).Interface()
		// 完全相同
		if (reflect.DeepEqual(itemValue, elementValue)) {
			index = i
			break
		}
	}

	return index
}

func Md5String(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)

	return hex.EncodeToString(hash.Sum(nil))
}

// IPv4 转换成数字
func IP4ToInt(ip string) int {
	nums := strings.Split(ip, ".")
	sum := 0
	for i := 0; i < len(nums); i++ {
		n, _ := strconv.Atoi(nums[i])
		sum += n
		sum <<= 8
	}

	return sum >> 8
}