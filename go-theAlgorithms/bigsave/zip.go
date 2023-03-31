package bigsave

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
)

// 实现了多种压缩算法的解压和压缩功能
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func GZipBytes(data []byte) []byte {
	var input bytes.Buffer
	g := gzip.NewWriter(&input)
	g.Write(data)
	g.Close()
	return input.Bytes()
}

func UGZipBytes(data []byte) []byte {
	var out bytes.Buffer //直接定义一个Buffer变量，不用初始化，可以直接使用
	var in bytes.Buffer
	in.Write(data)
	r, _ := gzip.NewReader(&in)
	r.Close()
	io.Copy(&out, r)
	return out.Bytes()

}
func ZipBytes(data []byte) []byte {

	var in bytes.Buffer
	z := zlib.NewWriter(&in)
	z.Write(data)
	z.Close()
	return in.Bytes()
}

func UZipBytes(data []byte) []byte {
	var out bytes.Buffer
	var in bytes.Buffer
	in.Write(data)
	r, _ := zlib.NewReader(&in)
	r.Close()
	io.Copy(&out, r)
	return out.Bytes()
}
