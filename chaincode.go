
/*
For iNSIDE TRACK code fest event (May-06 May-07 , 2017)
IBM Team : Trail Blazers
TeamID: TEAM17041405161300361
IdeaID:  IDEA17042017182300110
Summary: Chaincode ~ Smartcontract for Telco's to settle CDR's using Blockchain Smartcontracts;Eleminating Data Clearing Houses(DCH). 
The Smartcontract will reside on blockchain where Telco's will be member nodes and will agree to following Smartcontract.
All CDR processed will be preserved in blockchain as LedgerEntry and telco's can settle credit/debit based on txn summary and cdr reconsilation .
==Blockchain is TRUST building entity here replacing DCH
==Smartcontract will serve as processing and invoicing engine to share and settle CDR,Invoices amongst Telco Roaming partners in near realtime.
==Saving Time , Money for all member Telcos
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//LedgerEntry - telco's can store transaction in Ledger in this structure
type LedgerEntry struct {

	TxnType string `json:"txntype"`
	TxnAmount float64 `json:"txnamount"`
	InvoiceTo string `json:"invoiceto"`
	InvoiceFrom string `json:"invoiceto"`
	CDRid string `json:"cdrid"`

}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init 
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("init_chaincode", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// invoke different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "writetxnsummary" {
		return t.writeTxnSummay(stub, args)
	} else if function == "addledgerentry" {
		return t.addLedgerEntry(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "readtxntotal" { //read a variable
		return t.readTxnSummary(stub, args)
	} else if function == "readledgerentry" {
		return t.readLedgerEntryByCDRid(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) addLedgerEntry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding transaction as addLedgerEntry in blockchain")
	if len(args) != 5 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 5 for addLedgerEntry")
	}
	amt, err := strconv.ParseFloat(args[1], 64)
	

	ledgerentry := LedgerEntry{
		TxnType:   args[0],
		TxnAmount: amt,
		InvoiceTo: args[2],
		InvoiceFrom: args[3],
		CDRid: args[4],
	}

	bytes, err := json.Marshal(ledgerentry)
	if err != nil {
		fmt.Println("Error marshaling ledgerentry")
		return nil, errors.New("Error marshaling ledgerentry")
	}

	err = stub.PutState(ledgerentry.CDRid, bytes)
	if err != nil {
		return nil, err
}
return nil, nil
}


func (t *SimpleChaincode) readLedgerEntryByCDRid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("start readLedgerEntryByCDRid")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // cdrid
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	return bytes, nil
}


func (t *SimpleChaincode) writeTxnSummay(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] 
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) readTxnSummary(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}



    
