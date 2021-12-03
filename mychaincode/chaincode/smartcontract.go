package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract object provides functions for managing an Asset

type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset it a  business object
// ID is unique identifier
// Owner can be any one-- prosumer, consumer, utility-grid operator

type Asset struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Energy         string `json:"energy"`
	Price          int    `json:"price"`
	tType          string `json:"type"`
}


type Counter struct {
	_tokenId int64
}

// This is global variable
var counter Counter

// Increase the tokenId by one
func TokenIdIncrement() {
	counter._tokenId += 1
}

// Get the current tokenId
func TokenIdCurrent() int64 {
	return counter._tokenId
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, Owner string, Energy string, Price int, tType string) (string, error) {
	TokenIdIncrement()
	var ID = strconv.FormatInt(TokenIdCurrent(), 16)
	asset := Asset{
		ID:             ID,
		Owner:          Owner,
		Energy:         Energy,
		Price:       	Price,
		tType:          tType,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "-1", err
	}

	return ID, ctx.GetStub().PutState(ID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ID string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", ID)
	}
	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
