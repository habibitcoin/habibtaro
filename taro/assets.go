package taro

import (
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

type TaroExportProofRequest struct {
	AssetId   string `json:"asset_id"`
	ScriptKey string `json:"script_key"`
}

type TaroDecodeProofRequest struct {
	RawProof string `json:"raw_proof"`
	// proofindex
}

type TaroProofResponse struct {
	Proof TaroProof `json:"decoded_proof"`
}

type TaroProof struct {
	ProofIndex      string            `json:"proof_index"`
	NumberOfProofs  string            `json:"number_of_proofs"`
	Asset           TaroAssetResponse `json:"asset"`
	TxMerkleProof   string            `json:"tx_merkle_proof"`
	InclusionProof  string            `json:"inclusion_proof"`
	ExclusionProofs []string          `json:"exclusion_proofs"`
}

func (client *TaroClient) GetAssetProof(assetName string) (proofResponse TaroProofResponse, err error) {
	var assets TaroAssetsResponse
	var asset TaroAssetResponse
	if len(client.CachedAssetResponse.TaroAssets) == 0 {
		assets, err = client.ListAssets()
		if err != nil {
			return proofResponse, nil
		}
	} else {
		assets = client.CachedAssetResponse
	}

	for _, a := range assets.TaroAssets {
		if a.AssetGenesis.Name == assetName {
			asset = a
			break
		}
	}

	var encodedProof TaroDecodeProofRequest
	resp, err := client.sendPostRequestJSON("v1/taro/proofs/export", &TaroExportProofRequest{
		asset.AssetGenesis.AssetID,
		asset.ScriptKey,
	})
	if err != nil {
		log.Println(err)
		return proofResponse, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return proofResponse, err
	}

	if err := json.Unmarshal(bodyBytes, &encodedProof); err != nil {
		log.Println(err)
		return proofResponse, err
	}

	resp, err = client.sendPostRequestJSON("v1/taro/proofs/decode", &TaroDecodeProofRequest{
		encodedProof.RawProof,
	})
	if err != nil {
		log.Println(err)
		return proofResponse, err
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return proofResponse, err
	}

	if err := json.Unmarshal(bodyBytes, &proofResponse); err != nil {
		log.Println(err)
		return proofResponse, err
	}

	return proofResponse, err
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

	client.CachedAssetResponse = assets

	return assets, err
}
