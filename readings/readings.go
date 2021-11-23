package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Reading struct {
	ID       string `json:"ID"`
	Cycle    int    `json:"cycle"`
	Consumed int    `json:"consumed"`
	Injected int    `json:"injected"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	readings := []Reading{
		{ID: "reading1", Cycle: 1, Consumed: 5, Injected: 10},
		{ID: "reading2", Cycle: 1, Consumed: 20, Injected: 5},
	}

	for _, reading := range readings {
		readingJSON, err := json.Marshal(reading)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(reading.ID, readingJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateReading(ctx contractapi.TransactionContextInterface, id string, cycle int, consumed int, injected int) error {
	exists, err := s.ReadingExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the reading %s already exists", id)
	}

	reading := Reading{
		ID:       id,
		Cycle:    cycle,
		Consumed: consumed,
		Injected: injected,
	}
	readingJSON, err := json.Marshal(reading)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, readingJSON)
}

func (s *SmartContract) ReadReading(ctx contractapi.TransactionContextInterface, id string) (*Reading, error) {
	readingJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if readingJSON == nil {
		return nil, fmt.Errorf("the reading %s does not exist", id)
	}

	var reading Reading
	err = json.Unmarshal(readingJSON, &reading)
	if err != nil {
		return nil, err
	}

	return &reading, nil
}

func (s *SmartContract) GetAllReadings(ctx contractapi.TransactionContextInterface) ([]*Reading, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var readings []*Reading
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var reading Reading
		err = json.Unmarshal(queryResponse.Value, &reading)
		if err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}

func (s *SmartContract) ReadingExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	readingJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return readingJSON != nil, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating reading-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting reading-transfer-basic chaincode: %v", err)
	}
}
