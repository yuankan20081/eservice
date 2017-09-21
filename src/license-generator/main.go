package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	licenseCiper "game-util/license"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type License struct {
	Dir    string
	Server string
	BindIP string
	Expire time.Time
}

type PreperLicense struct {
	Server      string
	BindIP      string
	ExpireHours int64
}

func main() {
	descFilePath, _ := filepath.Abs("./description.json")
	prepare, err := ioutil.ReadFile(descFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var descriptions []*PreperLicense

	if err = json.Unmarshal(prepare, &descriptions); err != nil {
		log.Fatalln(err)
	}

	for _, desc := range descriptions {
		var l = License{
			Dir:    "some directory info goes here when loading to memory",
			Server: desc.Server,
			Expire: time.Now().Add(time.Hour * time.Duration(desc.ExpireHours)),
			BindIP: "127.0.0.1",
		}

		bin, err := json.Marshal(&l)
		if err != nil {
			log.Fatalln(err)
		}

		encrypted, err := licenseCiper.RsaEncrypt(bin)
		if err != nil {
			log.Fatalln(err)
		}

		txt := base64.StdEncoding.EncodeToString(encrypted)

		barr := md5.Sum([]byte(txt))
		fileName := fmt.Sprintf("./%s-%s.license", hex.EncodeToString(barr[:]), l.Server)

		if err := ioutil.WriteFile(fileName, []byte(txt), 0); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("license file generated:\n\tFile - %s\n\tExpires at - %s\n", filepath.Base(fileName), l.Expire)
		fmt.Println()
		fmt.Println()
	}
}
