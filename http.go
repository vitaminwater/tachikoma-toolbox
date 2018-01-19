package tachikoma

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var client *http.Client

func init() {
	proxyUser := os.Getenv("PROXY_USER")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	if proxyUser == "" && proxyPassword == "" {
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", proxyUser, proxyPassword)))

	ph := http.Header{}
	ph.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", auth))
	ph.Set("X-BOTPROXY-SESSION", strings.Replace(uuid.New().String(), "-", "", -1))
	log.Info(ph.Get("X-BOTPROXY-SESSION"))

	pu, err := url.Parse("http://x.botproxy.net:8080")
	Fatal(err)

	tr := &http.Transport{
		Proxy:              http.ProxyURL(pu),
		ProxyConnectHeader: ph,
	}
	client = &http.Client{Transport: tr}
}

func GetJSONProxy(url string, o interface{}) {
	r, err := client.Get(url)
	Fatal(err)
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(o)
	Fatal(err)
}

func GetJSONDirect(url string, o interface{}) {
	r, err := http.Get(url)
	Fatal(err)
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(o)
	Fatal(err)
}

func GetJSON(url string, o interface{}) {
	if client != nil {
		GetJSONProxy(url, o)
		return
	}
	GetJSONDirect(url, o)
}
