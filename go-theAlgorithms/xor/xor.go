package xor

/*
Encrypt使用Xor加密
将每个字节与key进行异或计算,返回值可能不可读，因为没有保证
在ASCII范围内
如果使用其他类型，如string， []int或其他类型，
添加将类型转换为[]byte的语句
*/
func Encrypt(key byte, plaintext []byte) []byte {
	cipherText := []byte{}
	for _, ch := range plaintext {
		cipherText = append(cipherText, key^ch)
	}
	return cipherText
}

/*
Decrypt使用Xor加密进行解密
*/
func Decrypt(key byte, cipherText []byte) []byte {
	plainText := []byte{}
	for _, ch := range cipherText {
		plainText = append(plainText, key^ch)
	}
	return plainText
}
