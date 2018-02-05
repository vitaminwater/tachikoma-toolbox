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
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

/**
 * error
 */

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * Http proxy client
 */

var client *http.Client

func init() {
	proxyUser := os.Getenv("PROXY_USER")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	if proxyUser == "" && proxyPassword == "" {
		client = http.DefaultClient
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

func GetTransport() http.RoundTripper {
	return client.Transport
}

/**
 * JSON
 */

func GetJSONDirect(url string, o interface{}) {
	r, err := http.Get(url)
	Fatal(err)
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(o)
	Fatal(err)
}

func GetJSON(url string, o interface{}) {
	r, err := client.Get(url)
	Fatal(err)
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(o)
	Fatal(err)
}

/**
 * Unmap
 */

func Unmap(v map[string]interface{}, t interface{}) {
	config := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: t}

	decoder, err := mapstructure.NewDecoder(&config)
	Fatal(err)

	err = decoder.Decode(v)
	Fatal(err)
}

/**
 * Float
 */

type Float float64

func (f *Float) UnmarshalJSON(b []byte) error {
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	fl, err := n.Float64()
	*f = Float(fl)
	return err
}

func (f Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}
