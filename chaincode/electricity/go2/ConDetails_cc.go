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
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the consumer personal and connection details structure, with 8 properties.  Structure tags are used by encoding/json library
type Con struct {
	Name   string `json:"name"`
	Address  string `json:"address"`
	MobileNo  string `json:"mobile_no"`
	AadharNo  string `json:"aadhar_no"`
	MeterNo  string `json:"meter_no"`
	Load  string `json:"load"`
	TarrifType  string `json:"tarrif_type"`
	AreaType  string `json:"area_type"`
}

/*
 * The Init method is called when the chaincode "ConDetails_cc" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the chaincode "ConDetails_cc"
 * The calling application program has also specified the particular chaincode function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested chaincode function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCon" {
		return s.queryCon(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "newCon" {
		return s.newCon(APIstub, args)
	} else if function == "queryAllCons" {
		return s.queryAllCons(APIstub)
	} else if function == "changeOwnership" {
		return s.changeOwnership(APIstub, args)
	} else if function == "changeMob" {
		return s.changeMob(APIstub, args)
	} else if function == "updateAdd" {
		return s.updateAdd(APIstub, args)
	} else if function == "updateAadhar" {
		return s.updateAadhar(APIstub, args)
	} else if function == "changeMeter" {
		return s.changeMeter(APIstub, args)
	} else if function == "changeLoad" {
		return s.changeLoad(APIstub, args)
	} else if function == "changeTarrifType" {
		return s.changeTarrifType(APIstub, args)
	} else if function == "changeAreaType" {
		return s.changeAreaType(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//initialising ledger with some consumer details
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cons := []Con{
		Con{Name: "Shirish", Address: "Purandarpur", MobileNo: "6202533516", AadharNo:"4683389079", MeterNo: "a435", Load: "1.25", TarrifType:"domestic", AreaType:"urban"},
		Con{Name: "Goku", Address: "Kankarbagh", MobileNo: "6202343516", AadharNo:"4683769079", MeterNo: "b783", Load: "1.50", TarrifType:"non-domestic", AreaType:"urban"},
		Con{Name: "Naruto", Address: "Boring road", MobileNo: "6209633516", AadharNo:"4639389079", MeterNo: "k554", Load: "1.75", TarrifType:"agriculture", AreaType:"urban"},
		Con{Name: "Sasuke", Address: "Patliputra", MobileNo: "6219533516", AadharNo:"4685589079", MeterNo: "e381", Load: "2.0", TarrifType:"industrial", AreaType:"urban"},
		}

	i := 0
	for i < len(cons) {
		fmt.Println("i is ", i)
		conAsBytes, _ := json.Marshal(cons[i])
		APIstub.PutState("CON"+strconv.Itoa(i), conAsBytes)
		fmt.Println("Added", cons[i])
		i = i + 1
	}

	return shim.Success(nil)
}
 
//view a particular consumer details
func (s *SmartContract) queryCon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(conAsBytes)
}

//new connection
func (s *SmartContract) newCon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	var con = Con{Name: args[1], Address: args[2], MobileNo: args[3], AadharNo: args[4], MeterNo: args[5], Load: args[6], TarrifType: args[7], AreaType: args[8]}

	conAsBytes, _ := json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//query all consumer details
func (s *SmartContract) queryAllCons(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CON0"
	endKey := "CON999"

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

	fmt.Printf("- queryAllCons:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//change consumer ownership
func (s *SmartContract) changeOwnership(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.Name = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer mobile
func (s *SmartContract) changeMob(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.MobileNo = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer address
func (s *SmartContract) updateAdd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.Address = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer aadhar no
func (s *SmartContract) updateAadhar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.AadharNo = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer meter no
func (s *SmartContract) changeMeter(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.MeterNo = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer load
func (s *SmartContract) changeLoad(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.Load = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer tarrif type
func (s *SmartContract) changeTarrifType(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.TarrifType = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

	return shim.Success(nil)
}

//change consumer area type
func (s *SmartContract) changeAreaType(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	conAsBytes, _ := APIstub.GetState(args[0])
	con := Con{}

	json.Unmarshal(conAsBytes, &con)
	con.AreaType = args[1]

	conAsBytes, _ = json.Marshal(con)
	APIstub.PutState(args[0], conAsBytes)

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

