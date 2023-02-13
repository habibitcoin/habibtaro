package taro

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
)

type TaroAssetsResponse struct {
	TaroAssets []TaroAssetResponse `json:"assets"`
}

type TaroAssetResponse struct {
	Version      int `json:"version"`
	AssetGenesis struct {
		GenesisPoint         string `json:"genesis_point"`
		Name                 string `json:"name"`
		Meta                 string `json:"meta"`
		AssetID              string `json:"asset_id"`
		OutputIndex          int    `json:"output_index"`
		GenesisBootstrapInfo string `json:"genesis_bootstrap_info"`
		Version              int    `json:"version"`
	} `json:"asset_genesis"`
	AssetType        string `json:"asset_type"`
	Amount           string `json:"amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	ScriptVersion    int    `json:"script_version"`
	ScriptKey        string `json:"script_key"`
	AssetGroup       string `json:"asset_group"`
	ChainAnchor      struct {
		AnchorTx        string `json:"anchor_tx"`
		AnchorTxid      string `json:"anchor_txid"`
		AnchorBlockHash string `json:"anchor_block_hash"`
		AnchorOutpoint  string `json:"anchor_outpoint"`
		InternalKey     string `json:"internal_key"`
	} `json:"chain_anchor"`
	PrevWitnesses []interface{} `json:"prev_witnesses"`
}

func (client *TaroClient) GetAsset(assetName string) (assetResponse TaroAssetResponse, err error) {
	var assets TaroAssetsResponse
	if len(client.CachedAssetResponse.TaroAssets) == 0 {
		assets, err = client.ListAssets()
		if err != nil {
			return assetResponse, nil
		}
	} else {
		assets = client.CachedAssetResponse
	}

	for _, asset := range assets.TaroAssets {
		if asset.AssetGenesis.Name == assetName {
			assetResponse = asset
			break
		}
	}

	return assetResponse, err
}

func (client *TaroClient) ListAssets() (assets TaroAssetsResponse, err error) {
	resp, err := client.sendGetRequest("v1/taro/assets")
	if err != nil {
		log.Println(err)
		return assets, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return assets, err
	}

	if err := json.Unmarshal(bodyBytes, &assets); err != nil {
		log.Println(err)
		return assets, err
	}

	for i, a := range assets.TaroAssets {
		str, _ := base64.StdEncoding.DecodeString(a.AssetGenesis.GenesisBootstrapInfo)
		assets.TaroAssets[i].AssetGenesis.GenesisBootstrapInfo = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.AssetGenesis.Meta)
		assets.TaroAssets[i].AssetGenesis.Meta = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.AssetGenesis.GenesisBootstrapInfo)
		assets.TaroAssets[i].AssetGenesis.AssetID = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ScriptKey)
		assets.TaroAssets[i].ScriptKey = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ChainAnchor.AnchorTx)
		assets.TaroAssets[i].ChainAnchor.AnchorTx = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ChainAnchor.InternalKey)
		assets.TaroAssets[i].ChainAnchor.InternalKey = hex.EncodeToString(str)

	}

	client.CachedAssetResponse = assets

	return assets, err
}
