package taro

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/habibitcoin/habibtaro/configs"
)

type TaroClient struct {
	Client   *http.Client
	Host     string
	Macaroon string
	Context  context.Context
}

// func NewTaroClient
func NewClient(ctx context.Context) (client TaroClient) {
	config := configs.GetConfig(ctx)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
	}
	client = TaroClient{
		Client:   httpClient,
		Host:     config.TaroHost,
		Macaroon: loadMacaroon(ctx),
	}

	return client
}

func loadMacaroon(ctx context.Context) (macaroon string) {
	macaroonBytes, err := ioutil.ReadFile(configs.GetConfig(ctx).TaroMacaroonLocation)
	if err != nil {
		log.Println("couldnt find or open macaroon")
		log.Println(err)
		return configs.GetConfig(ctx).TaroMacaroon
	}

	macaroon = hex.EncodeToString(macaroonBytes)

	log.Println(macaroon)
	return macaroon
}

func (client *TaroClient) sendGetRequest(endpoint string) (*http.Response, error) {
	log.Println(client.Host + endpoint)
	req, err := http.NewRequest("GET", client.Host+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Grpc-Metadata-macaroon", client.Macaroon)
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (client *TaroClient) sendPostRequestJSON(endpoint string, payload interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(payload)
	req, err := http.NewRequest("POST", client.Host+endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Grpc-Metadata-macaroon", client.Macaroon)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (client *TaroClient) sendPostRequest(endpoint string, payload string) (*http.Response, error) {
	jsonStr := []byte(payload)

	req, err := http.NewRequest("POST", client.Host+endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Grpc-Metadata-macaroon", client.Macaroon)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
