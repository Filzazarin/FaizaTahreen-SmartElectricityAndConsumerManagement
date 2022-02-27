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
 
 // Define the payment details structure, with 6 properties.  Structure tags are used by encoding/json library
 type Con struct {
	 Name   string `json:"name"`
	 BillNo  string `json:"bill_no"`
	 ReceiptNo string `json:"receipt_no"`
	 AmtPaid  string `json:"amount_paid"`
	 EnergyConsumption  string `json:"energy_consumption"`
	 PaymentDate  string `json:"payment_date"`
 }
 
 /*
  * The Init method is called when the chaincode "PayDetails_cc" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function -- see initLedger()
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the chaincode "PayDetails_cc"
  * The calling application program has also specified the particular chaincode function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested chaincode function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "viewPayDetails" {
		 return s.viewPayDetails(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "addPayment" {
		 return s.addPayment(APIstub, args)
	 } else if function == "viewAllPayDetails" {
		 return s.viewAllPayDetails(APIstub)
	 } else if function == "updateAmountPaid" {
		 return s.updateAmountPaid(APIstub, args)
	 } else if function == "updateEnergyConsumption" {
		 return s.updateEnergyConsumption(APIstub, args)
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 //initialising ledger with some consumer payment details
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 cons := []Con{
		 Con{Name: "Shirish", BillNo: "da-345", ReceiptNo: "rn-882", AmtPaid: "1654", EnergyConsumption:"85", PaymentDate:"15-08-2020"},
		 Con{Name: "Goku", BillNo: "ew-234", ReceiptNo: "iu-120", AmtPaid: "2887", EnergyConsumption:"172", PaymentDate:"12-08-2020"},
		 Con{Name: "Naruto", BillNo: "ws-287", ReceiptNo: "oc-365", AmtPaid: "889", EnergyConsumption:"65", PaymentDate:"16-08-2020"},
		 Con{Name: "Sasuke", BillNo: "kj-419", ReceiptNo: "sn-214", AmtPaid: "1256", EnergyConsumption:"112", PaymentDate:"17-08-2020"},
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
 
 //view a consumer's payment details
 func (s *SmartContract) viewPayDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 conAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(conAsBytes)
 }
 
 //add a consumer's payment details
 func (s *SmartContract) addPayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 7 {
		 return shim.Error("Incorrect number of arguments. Expecting 7")
	 }
 
	 var con = Con{Name: args[1], BillNo: args[2], ReceiptNo: args[3], AmtPaid: args[4], EnergyConsumption: args[5], PaymentDate: args[6]}
 
	 conAsBytes, _ := json.Marshal(con)
	 APIstub.PutState(args[0], conAsBytes)
 
	 return shim.Success(nil)
 }
 
 //view payment details of all consumers
 func (s *SmartContract) viewAllPayDetails(APIstub shim.ChaincodeStubInterface) sc.Response {
 
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
 
 //update the amount paid
 func (s *SmartContract) updateAmountPaid(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 conAsBytes, _ := APIstub.GetState(args[0])
	 con := Con{}
 
	 json.Unmarshal(conAsBytes, &con)
	 con.AmtPaid = args[1]
 
	 conAsBytes, _ = json.Marshal(con)
	 APIstub.PutState(args[0], conAsBytes)
 
	 return shim.Success(nil)
 }

 //update the energy consumption
 func (s *SmartContract) updateEnergyConsumption(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 conAsBytes, _ := APIstub.GetState(args[0])
	 con := Con{}
 
	 json.Unmarshal(conAsBytes, &con)
	 con.EnergyConsumption = args[1]
 
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
 


		
