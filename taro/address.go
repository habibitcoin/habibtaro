package taro

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
)

type TaroAddressRequest struct {
	BootstrapInfo string `json:"genesis_bootstrap_info" form:"genesis_bootstrap_info" bson:"genesis_bootstrap_info"`
	GroupKey      string `json:"group_key" form:"group_key" bson:"group_key"`
	Amount        string `json:"amt" form:"amt" bson:"amt"`
}

type TaroAddressesResponse struct {
	TaroAssets []TaroAssetResponse `json:"assets"`
}

type TaroAddressResponse struct {
	Address          string `json:"encoded"`
	AssetId          string `json:"asset_id"`
	AssetType        string `json:"asset_type"`
	Amount           string `json:"amount"`
	GroupKey         string `json:"group_key"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TaprootOutputKey string `json:"taproot_output_key"`
}

func (client *TaroClient) CreateAddress(genesisBootstrapInfo, groupKey, amt string) (address TaroAddressResponse, err error) {
	log.Println(genesisBootstrapInfo)
	bootstrapHex, err := hex.DecodeString(genesisBootstrapInfo)
	if err != nil {
		return address, err
	}
	genesisBootstrapInfo = base64.URLEncoding.EncodeToString(bootstrapHex)
	log.Println(genesisBootstrapInfo)

	if groupKey != "" {
		groupKeyHex, err := hex.DecodeString(groupKey)
		if err != nil {
			return address, err
		}
		groupKey = base64.URLEncoding.EncodeToString(groupKeyHex)
	}

	resp, err := client.sendPostRequestJSON("v1/taro/addrs", &TaroAddressRequest{
		genesisBootstrapInfo,
		groupKey,
		amt,
	})
	if err != nil {
		log.Println(err)
		return address, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return address, err
	}

	if err := json.Unmarshal(bodyBytes, &address); err != nil {
		log.Println(err)
		return address, err
	}

	str, _ := base64.StdEncoding.DecodeString(address.AssetId)
	address.AssetId = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.GroupKey)
	address.GroupKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.ScriptKey)
	address.ScriptKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.InternalKey)
	address.InternalKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.TaprootOutputKey)
	address.TaprootOutputKey = hex.EncodeToString(str)

	return address, err
}
