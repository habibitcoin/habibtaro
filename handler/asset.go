package handler

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/habibitcoin/habibtaro/configs"
	"github.com/habibitcoin/habibtaro/taro"
	"github.com/labstack/echo"
)

func (h *Handler) FetchAsset(c echo.Context) (err error) {
	var (
		ctx = context.Background()
	)

	ctx, _ = configs.LoadConfig(ctx)
	taroClient := taro.NewClient(ctx)

	resp, err := taroClient.GetAsset(c.Param("name"))
	if err != nil {
		log.Println("Error creating address")
		log.Println(err)
		return err
	}

	resp = decodeAssetFields(resp)

	str, _ := json.MarshalIndent(resp, "", "\t")

	str1 := strings.Replace(string(str), "\n", "<br>", -1)
	str1 = strings.Replace(string(str1), "\t", "&emsp;", -1)

	return c.Render(http.StatusOK, "asset.html", map[string]interface{}{
		"asset":     template.HTML(str1),
		"assetName": resp.AssetGenesis.Name,
	})
}

func (h *Handler) FetchAssetProof(c echo.Context) (err error) {
	var (
		ctx = context.Background()
	)

	ctx, _ = configs.LoadConfig(ctx)
	taroClient := taro.NewClient(ctx)

	resp, err := taroClient.GetAssetProof(c.Param("name"))
	if err != nil {
		log.Println("Error exporting proof")
		log.Println(err)
		return err
	}

	resp = decodeProofFields(resp)

	str, _ := json.MarshalIndent(resp, "", "\t")

	str1 := strings.Replace(string(str), "\n", "<br>", -1)
	str1 = strings.Replace(string(str1), "\t", "&emsp;", -1)

	return c.Render(http.StatusOK, "proof.html", map[string]interface{}{
		"proof": template.HTML(str1),
	})
}

func decodeProofFields(p taro.TaroProofResponse) (finalProof taro.TaroProofResponse) {
	str, _ := base64.StdEncoding.DecodeString(p.Proof.TxMerkleProof)
	p.Proof.TxMerkleProof = hex.EncodeToString(str)
	str, _ = base64.StdEncoding.DecodeString(p.Proof.InclusionProof)
	p.Proof.InclusionProof = hex.EncodeToString(str)
	for i, e := range p.Proof.ExclusionProofs {
		str, _ = base64.StdEncoding.DecodeString(e)
		p.Proof.ExclusionProofs[i] = hex.EncodeToString(str)
	}
	p.Proof.Asset = decodeAssetFields(p.Proof.Asset)

	return p
}

func decodeAssetFields(a taro.TaroAssetResponse) (finalAsset taro.TaroAssetResponse) {
	str, _ := base64.StdEncoding.DecodeString(a.AssetGenesis.GenesisBootstrapInfo)
	a.AssetGenesis.GenesisBootstrapInfo = hex.EncodeToString(str)
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

	return a
}
