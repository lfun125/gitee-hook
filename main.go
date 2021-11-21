package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gitee-hook/config"
	"gitee-hook/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	Config config.Model
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := do(rw, r); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(err.Error()))
		}
	})
	panic(http.ListenAndServe(Config.Listen, nil))
}

func do(rw http.ResponseWriter, r *http.Request) error {
	var data models.Data
	requestRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(bytes.NewReader(requestRaw)).Decode(&data); err != nil {
		return err
	}
	reqTimestamp := r.Header.Get("X-Gitee-Timestamp")
	if reqTimestamp == "" {
		return errors.New("Header not set X-Gitee-Timestamp")
	}
	if data.Sign != sign(reqTimestamp, Config.GiteeSecret) {
		return errors.New("Gitee sign error")
	}

	repository, ok := Config.Repository[data.Repository.FullName]
	if !ok {
		return errors.New("Not config this repository")
	}
	var job string
	for b, j := range repository.Branches {
		if path.Base(data.Ref) == b {
			job = j
			break
		}
	}
	if job == "" {
		rw.WriteHeader(http.StatusOK)
		echo(rw, `This branch does not need to be released`)
		return nil
	}
	requestUrl := fmt.Sprintf("%s/job/%s/build?token=%s", strings.TrimRight(Config.JenkinsUrl, "/"), job, Config.JenkinsProjectToken)

	req, err := http.NewRequest(r.Method, requestUrl, bytes.NewReader(requestRaw))
	if err != nil {
		return err
	}
	req.SetBasicAuth(Config.JenkinsUser, Config.JenkinsUserToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	rw.WriteHeader(resp.StatusCode)
	for k := range resp.Header {
		rw.Header().Set(k, resp.Header.Get(k))
	}
	if _, err := io.Copy(rw, resp.Body); err != nil {
		return err
	}
	return err
}

func sign(t, secret string) string {
	s := fmt.Sprintf("%s\n%s", t, secret)
	bts := hmacSha256(s, secret)
	v := base64.StdEncoding.EncodeToString(bts)
	return v
}

func hmacSha256(data string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return h.Sum(nil)
}

func init() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "f", "./config.yml", "config file")
	flag.Parse()
	configRaw, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if err := yaml.Unmarshal(configRaw, &Config); err != nil {
		log.Fatalln(err)
	}
}

func echo(w http.ResponseWriter, v interface{}) {
	if _, err := w.Write([]byte(fmt.Sprintf("%v", v))); err != nil {
		log.Println(err)
	}
}
