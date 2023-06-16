package taro

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TaroAddressRequest struct {
	AssetId string `json:"asset_id" form:"asset_id" bson:"asset_id"`
	Amount  string `json:"amt" form:"amt" bson:"amt"`
	// ScriptKey   string `json:"script_key" form:"script_key" bson:"script_key"`
	// InternalKey string `json:"internal_key" form:"internal_key" bson:"internal_key"`
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

func (client *TaroClient) CreateAddress(assetId, amt string) (address TaroAddressResponse, err error) {
	/*if groupKey != "" {
		groupKeyHex, err := hex.DecodeString(groupKey)
		if err != nil {
			return address, err
		}
		groupKey = base64.URLEncoding.EncodeToString(groupKeyHex)
	}*/

	resp, err := client.sendPostRequestJSON("v1/taproot-assets/addrs", &TaroAddressRequest{
		assetId, amt,
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

	/*str, _ := base64.StdEncoding.DecodeString(address.AssetId)
	address.AssetId = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.GroupKey)
	address.GroupKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.ScriptKey)
	address.ScriptKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.InternalKey)
	address.InternalKey = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(address.TaprootOutputKey)
	address.TaprootOutputKey = hex.EncodeToString(str)*/

	return address, err
}
