package license

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	licenseCiper "game-util/license"
	"io/ioutil"
	"testing"
	"time"
)

func TestWriteLiceseToFile(t *testing.T) {
	var l = License{
		Dir:    "useless",
		Server: "这是一个猜点测试服",
		Expire: time.Now().Add(time.Hour * 1),
	}

	bin, err := json.Marshal(&l)
	if err != nil {
		t.Error(err)
	}

	encrypted, err := licenseCiper.RsaEncrypt(bin)
	if err != nil {
		t.Error(err)
	}

	txt := base64.StdEncoding.EncodeToString(encrypted)

	barr := md5.Sum([]byte(txt))
	fileName := fmt.Sprintf("C:/Users/yuank/Desktop/%s.license", hex.EncodeToString(barr[:]))

	if err := ioutil.WriteFile(fileName, []byte(txt), 0); err != nil {
		t.Error(err)
	}
}
