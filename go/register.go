package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//Seller struct
type Seller struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`
	Dimension  string `json:"dimension"`
	Locality   string `json:"locality"`
	Landprice  string `json:"landprice"`
	Owner      string `json:"owner"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	// Handle different functions
	switch function {
	case "registerLand":
		return t.registerLand(stub, args)
	case "fetchLand":
		return t.fetchLandDetails(stub, args)
	case "transferLand":
		return t.transferLand(stub, args)
	default:
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

func (t *SimpleChaincode) registerLand(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var seller = Seller{Name: args[0], Dimension: args[1], Locality: args[2], Landprice: args[3], Owner: args[4]}

	sellerAsBytes, _ := json.Marshal(seller)
	APIstub.PutState(args[0], sellerAsBytes)

	return shim.Success(nil)
}

func (t *SimpleChaincode) fetchLandDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the land to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"land does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *SimpleChaincode) transferLand(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	sellerName := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferland ", sellerName, newOwner)

	sellerAsBytes, err := stub.GetState(sellerName)
	if err != nil {
		return shim.Error("Failed to get land:" + err.Error())
	} else if sellerAsBytes == nil {
		return shim.Error("land does not exist")
	}

	landToTransfer := Seller{}
	err = json.Unmarshal(sellerAsBytes, &landToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	landToTransfer.Owner = newOwner //change the owner

	landJSONasBytes, _ := json.Marshal(landToTransfer)
	err = stub.PutState(sellerName, landJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferland (success)")
	return shim.Success(nil)
}
