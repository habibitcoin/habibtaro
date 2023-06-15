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
		GenesisPoint string `json:"genesis_point"`
		Name         string `json:"name"`
		Meta         string `json:"meta_hash"`
		AssetID      string `json:"asset_id"`
		OutputIndex  int    `json:"output_index"`
		Version      int    `json:"version"`
	} `json:"asset_genesis"`
	AssetType        string `json:"asset_type"`
	Amount           string `json:"amount"`
	LockTime         int    `json:"lock_time"`
	RelativeLockTime int    `json:"relative_lock_time"`
	ScriptVersion    int    `json:"script_version"`
	ScriptKey        string `json:"script_key"`
	ScriptKeyIsLocal bool   `json:"script_key_is_local"`
	AssetGroup       string `json:"asset_group"`
	ChainAnchor      struct {
		AnchorTx         string `json:"anchor_tx"`
		AnchorTxid       string `json:"anchor_txid"`
		AnchorBlockHash  string `json:"anchor_block_hash"`
		AnchorOutpoint   string `json:"anchor_outpoint"`
		InternalKey      string `json:"internal_key"`
		MerkleRoot       string `json:"merkle_root"`
		TapscriptSibling string `json:"tapscript_sibling"`
	} `json:"chain_anchor"`
	PrevWitnesses []struct {
		PrevId struct {
			AnchorPoint string `json:"anchor_point"`
			AssetId     string `json:"asset_id"`
			ScriptKey   string `json:"script_key"`
			Amount      string `json:"amount"`
		} `json:"prev_id"`
		TxWitness       []string `json:"tx_witness"`
		SplitCommitment string   `json:"split_commitment"`
	} `json:"prev_witnesses"`
	IsSpent bool `json:"is_spent"`
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
	RawProof          string `json:"raw_proof"`
	WithPrevWitnesses bool   `json:"with_prev_witnesses"`
	WithMetaReveal    bool   `json:"with_meta_reveal"`
	// proofindex
}

type TaroProofResponse struct {
	Proof TaroProof `json:"decoded_proof"`
}

type TaroProof struct {
	ProofAtDepth   string            `json:"proof_at_depth"`
	NumberOfProofs string            `json:"number_of_proofs"`
	Asset          TaroAssetResponse `json:"asset"`
	MetaReveal     struct {
		Data     string `json:"data"`
		Type     string `json:"type"`
		MetaHash string `json:"meta_hash"`
	} `json:"meta_reveal"`
	TxMerkleProof   string   `json:"tx_merkle_proof"`
	InclusionProof  string   `json:"inclusion_proof"`
	ExclusionProofs []string `json:"exclusion_proofs"`
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

	log.Println("hello2")

	var encodedProof TaroDecodeProofRequest
	resp, err := client.sendPostRequestJSON("v1/taproot-assets/proofs/export", &TaroExportProofRequest{
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

	resp, err = client.sendPostRequestJSON("v1/taproot-assets/proofs/decode", &TaroDecodeProofRequest{
		encodedProof.RawProof, true, true,
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

	log.Println(string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &proofResponse); err != nil {
		log.Println(err)
		return proofResponse, err
	}

	log.Println(proofResponse)

	return proofResponse, err
}

func (client *TaroClient) ListAssets() (assets TaroAssetsResponse, err error) {
	resp, err := client.sendGetRequest("v1/taproot-assets/assets")
	if err != nil {
		log.Println(err)
		return assets, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("death")
		log.Println(err)
		return assets, err
	}

	if err := json.Unmarshal(bodyBytes, &assets); err != nil {
		log.Println("death2")
		log.Println(err)
		return assets, err
	}

	client.CachedAssetResponse = assets

	return assets, err
}
