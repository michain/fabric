package houseregistry

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)


type House struct {
	ID	string `json:"id"`
	Owner  string `json:"owner"`
	Price  int `json:"price"`
	Colour string `json:"colour"`
}

// HouseRegistryChaincode House Registry Chaincode implementation
type HouseRegistryChaincode struct {
}

func (t *HouseRegistryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	for i:=0;i< 6;i++ {
		stub.PutState(strconv.Itoa(i), []byte("house-info-" + strconv.Itoa(i)))
	}

	return shim.Success(nil)
}

// queryHouse query house info with house id
func (t *HouseRegistryChaincode) queryHouse(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := stub.GetState(args[0])
	return shim.Success(carAsBytes)
}

// createHouse create house info with house json string
func (t *HouseRegistryChaincode) createHouse(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	stub.PutState(args[0], []byte(args[1]))

	return shim.Success(nil)
}


func (t *HouseRegistryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "queryHouse" {
		return t.queryHouse(stub, args)
	} else if function == "createHouse" {
		return t.createHouse(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"queryHouse\" \"createHouse\"")
}
