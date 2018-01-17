/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

/*
 * peer chaincode install -p chaincodedev/chaincode/cshbond -n mycc -v 0
 * peer chaincode instantiate -n mycc -v 0 -c '{"Args":[""]}' -C myc
 * peer chaincode invoke -n mycc -c '{"function":"initLedger","Args":[""]}' -C myc
 * peer chaincode query -n mycc -c '{"function":"queryAllBonds","Args":[""]}' -C myc
 * peer chaincode query -n mycc -c '{"function":"queryBond","Args":["BOND4"]}' -C myc
 * peer chaincode invoke -n mycc -c '{"function":"createBond","Args":["BOND11","Bond11","CS11","6","201612","202012"]}' -C myc
 * peer chaincode invoke -n mycc -c '{"function":"createBond","Args":["BOND10","Bond10","CS11","6","201612","202012"]}' -C myc
 * peer chaincode query -n mycc -c '{"function":"queryBond","Args":["BOND11"]}' -C myc
 * peer chaincode query -n mycc -c '{"function":"queryBond","Args":["BOND10"]}' -C myc
 */
package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the bond structure, with 5 properties.  Structure tags are used by encoding/json library
type Bond struct {
	BondID    string `json:"bondid"`
	Owner     string `json:"owner"`
	Rate      string `json:"rate"`
	IssueDate string `json:"issuedate"`
	DueDate   string `json:"duedate"`
}

//issuerId: 'issuer0',
//id: 'issuer0.2017.6.13.600',
//principal: 500000,
//term: 12,
//maturityDate: '2017.6.13',
//rate: 600,
//trigger: 'hurricane 2 FL',
//state: 'offer'

//type Bond struct {
//	IssuerId       string `json:"issuerId"`
//	Id             string `json:"id"`
//	Principal      uint64 `json:"principal"`
//	Term           uint64 `json:"term"`
//	MaturityDate   string `json:"maturityDate"`
//	Rate           uint64 `json:"rate"`
//	Trigger        string `json:"trigger"`
//	State          string `json:"state"`
//}

// Define the bond structure, with 6 properties.  Structure tags are used by encoding/json library
/*		0
	json
  	{
		"ticker":  "string",
		"par": 0.00,
		"qty": 10,
		"discount": 7.5,
		"maturity": 30,
		"owners": [ // This one is not required
			{
				"company": "company1",
				"quantity": 5
			},
			{
				"company": "company3",
				"quantity": 3
			},
			{
				"company": "company4",
				"quantity": 2
			}
		],
		"issuer":"company2",
		"issueDate":"1456161763790"  (current time in milliseconds as a string)
	}
*/
//type CP struct {
//	CUSIP     string  `json:"cusip"`
//	Ticker    string  `json:"ticker"`
//	Par       float64 `json:"par"`
//	Qty       int     `json:"qty"`
//	Discount  float64 `json:"discount"`
//	Maturity  int     `json:"maturity"`
//	Owners    []Owner `json:"owner"`
//	Issuer    string  `json:"issuer"`
//	IssueDate string  `json:"issueDate"`
//}

/*
 * The Init method is called when the Smart Contract "cshbond" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "cshbond"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryBond" {
		return s.queryBond(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createBond" {
		return s.createBond(APIstub, args)
	} else if function == "queryAllBonds" {
		return s.queryAllBonds(APIstub)
	} else if function == "changeBondOwner" {
		return s.changeBondOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryBond(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	bondAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(bondAsBytes)
}

//type Bond struct {
//	BondID     string `json:"bondid"`
//	Owner   string `json:"owner"`
//	Rate  float64 `json:"rate"`
//	IssueDate string `json:"issuedate"`
//	DueDate  string `json:"duedate"`
//}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	bonds := []Bond{
		Bond{BondID: "Bond1", Owner: "CS1", Rate: "1", IssueDate: "201601", DueDate: "202001"},
		Bond{BondID: "Bond2", Owner: "CS2", Rate: "1.5", IssueDate: "201602", DueDate: "202002"},
		Bond{BondID: "Bond3", Owner: "CS3", Rate: "2", IssueDate: "201603", DueDate: "202003"},
		Bond{BondID: "Bond4", Owner: "CS4", Rate: "2.5", IssueDate: "201604", DueDate: "202004"},
		Bond{BondID: "Bond5", Owner: "CS5", Rate: "3", IssueDate: "201605", DueDate: "202005"},
		Bond{BondID: "Bond6", Owner: "CS6", Rate: "3.5", IssueDate: "201607", DueDate: "202006"},
		Bond{BondID: "Bond7", Owner: "CS7", Rate: "4", IssueDate: "201608", DueDate: "202007"},
		Bond{BondID: "Bond8", Owner: "CS8", Rate: "4.5", IssueDate: "201609", DueDate: "202008"},
		Bond{BondID: "Bond9", Owner: "CS9", Rate: "5", IssueDate: "201610", DueDate: "202009"},
		Bond{BondID: "Bond10", Owner: "CS10", Rate: "5.5", IssueDate: "201611", DueDate: "202010"},
	}

	i := 0
	for i < len(bonds) {
		fmt.Println("i is ", i)
		bondAsBytes, _ := json.Marshal(bonds[i])
		APIstub.PutState("BOND"+strconv.Itoa(i), bondAsBytes)
		fmt.Println("Added", bonds[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createBond(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	//if len(args) != 5 {
	//	return shim.Error("Incorrect number of arguments. Expecting 5")
	//}
	//var bond = Bond{Name: args[1], Model: args[2], Colour: args[3], Owner: args[4]}
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	var bond = Bond{BondID: args[1], Owner: args[2], Rate: args[3], IssueDate: args[4], DueDate: args[5]}

	bondAsBytes, _ := json.Marshal(bond)
	APIstub.PutState(args[0], bondAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllBonds(APIstub shim.ChaincodeStubInterface) peer.Response {

	startKey := "BOND0"
	endKey := "BOND999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBonds:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeBondOwner(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	bondAsBytes, _ := APIstub.GetState(args[0])
	bond := Bond{}

	json.Unmarshal(bondAsBytes, &bond)
	bond.Owner = args[1]

	bondAsBytes, _ = json.Marshal(bond)
	APIstub.PutState(args[0], bondAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
