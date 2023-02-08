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
		log.Println(hex.EncodeToString([]byte(a.AssetGenesis.GenesisBootstrapInfo)))
		str, _ := base64.StdEncoding.DecodeString(a.AssetGenesis.GenesisBootstrapInfo)
		assets.TaroAssets[i].AssetGenesis.GenesisBootstrapInfo = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.AssetGenesis.Meta)
		a.AssetGenesis.Meta = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.AssetGenesis.GenesisBootstrapInfo)
		a.AssetGenesis.AssetID = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ScriptKey)
		a.ScriptKey = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ChainAnchor.AnchorTx)
		a.ChainAnchor.AnchorTx = hex.EncodeToString(str)
		str, _ = base64.StdEncoding.DecodeString(a.ChainAnchor.InternalKey)
		a.ChainAnchor.InternalKey = hex.EncodeToString(str)

	}

	return assets, err
}

/*func (a *TaroAssetResponse) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	a.AssetGenesis.GenesisBootstrapInfo = hex.EncodeToString([]byte(a.AssetGenesis.GenesisBootstrapInfo))
	a.AssetGenesis.Meta = hex.EncodeToString([]byte(a.AssetGenesis.Meta))
	a.AssetGenesis.AssetID = hex.EncodeToString([]byte(a.AssetGenesis.AssetID))
	a.ScriptKey = hex.EncodeToString([]byte(a.ScriptKey))
	a.ChainAnchor.AnchorTx = hex.EncodeToString([]byte(a.ChainAnchor.AnchorTx))
	a.ChainAnchor.InternalKey = hex.EncodeToString([]byte(a.ChainAnchor.InternalKey))

	return nil
}*/
