package encode

import (
	"crypto/md5"
	"fmt"
)

func EncodeMd5(s string) string {
	// 进行md5加密，因为Sum函数接受的是字节数组，因此需要注意类型转换
	srcCode := md5.Sum([]byte(s))

	// md5.Sum函数加密后返回的是字节数组，需要转换成16进制形式
	code := fmt.Sprintf("%x", srcCode)

	return string(code)
}
