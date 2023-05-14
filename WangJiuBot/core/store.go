package core

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"path"
	"reflect"
)

var (
	db      *bolt.DB
	Buckets = []Bucket{}
	WangJiu = NewBucket("WangJiu")
)

type Bucket string

func initStore() {
	var err error
	dir, err := os.UserConfigDir()
	confPath := path.Join(dir, "WangJiu.cache")
	db, err = bolt.Open(confPath, 0600, nil)
	//TODO:panic: open /root/.config/WangJiu.cache: no such file or directory 解决文件不存在出错
	if err != nil {
		panic(err)
	}
}
func (b Bucket) String() string {
	return string(b)
}
func NewBucket(name string) Bucket {
	b := Bucket(name)
	Buckets = append(Buckets, b)
	return b
}

func (b Bucket) Set(key, val interface{}) {
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		//如果没有这个数据库
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(b))
		}
		k := fmt.Sprint(key)
		v := fmt.Sprint(val)
		if v == "" {
			//如果val是空, 则删除key
			bucket.Delete([]byte(k))
		} else {
			bucket.Put([]byte(k), []byte(v))
		}

		return nil

	})
}

func (b Bucket) Get(key interface{}, defaultVal ...string) string {
	var k, value string
	if defaultVal != nil {
		if len(defaultVal) == 1 {
			value = defaultVal[0]
		}
	}
	k = fmt.Sprint(key)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket == nil {
			return nil
		}
		if v := string(bucket.Get([]byte(k))); v != "" {
			value = v
		}

		return nil
	})
	return value

}

func (b Bucket) GetInt(key interface{}, defaultVal ...int) int {
	var val int
	if len(defaultVal) != 0 {
		val = defaultVal[0]

	}
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket == nil {
			return nil
		}
		v := Int(string(bucket.Get([]byte(fmt.Sprint(key)))))
		if v != 0 {
			val = v
		}
		return nil
	})

	return val
}

func (b Bucket) GetBool(key interface{}, defaultVal ...bool) bool {
	var value bool
	if len(defaultVal) != 0 {
		value = defaultVal[0]
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket != nil {
			return nil
		}
		v := string(b.Get([]byte(fmt.Sprint(key))))
		switch v {
		case "true":
			value = true
		case "false":
			value = false
		}
		return nil

	})

	return value

}

func (b Bucket) Foreach(f func(k []byte, v []byte) error) {
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket != nil {
			bucket.ForEach(f)
		}
		return nil

	})
}

func (b Bucket) Create(i interface{}) error {
	s := reflect.ValueOf(i).Elem()
	id := s.FieldByName("ID")
	sequemce := s.FieldByName("Sequence")
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(b))
		}
		if _, ok := id.Interface().(int); ok {
			key := id.Int()
			sq, _ := bucket.NextSequence()
			if key == 0 {
				//说明原本属性ID是空值,使用生成的序列作为ID
				key = int64(sq)
				id.SetInt(key)
			}
			if sequemce != reflect.ValueOf(nil) {
				sequemce.SetInt(int64(sq))
			}
			buf, err := json.Marshal(i)
			if err != nil {
				return err
			}
			return bucket.Put(Itob(uint64(key)), buf)
		} else {
			key := id.String()
			sq, _ := bucket.NextSequence()
			if key == "" {
				key = fmt.Sprint(sq)
				id.SetString(key)
			}
			if sequemce != reflect.ValueOf(nil) {
				sequemce.SetInt(int64(sq))
			}
			buf, err := json.Marshal(i)
			if err != nil {
				return err
			}
			return bucket.Put([]byte(key), buf)

		}
	})
}
