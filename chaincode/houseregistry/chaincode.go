package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
	"errors"
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

	houses := []House{
		House{ID: "1", Owner: "张三", Price:100000, Colour: "blue"},
		House{ID: "2", Owner: "李四", Price:100000, Colour: "red"},
		House{ID: "3", Owner: "王五", Price:100000, Colour: "green"},
		House{ID: "4", Owner: "Tomoko", Price:100000, Colour: "yellow"},
		House{ID: "5", Owner: "Max", Price:100000, Colour: "black"},
		House{ID: "6", Owner: "Brad", Price:100000, Colour: "purple"},
	}


	for i:=0;i<len(houses);i++ {
		data, err := json.Marshal(houses[i])
		if err != nil{
			continue
		}
		stub.PutState(getHouseID(houses[i].ID), data)
	}

	return shim.Success(nil)
}

// queryHouse query house info with house id
func (t *HouseRegistryChaincode) queryHouse(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dataBytes, err := stub.GetState(getHouseID(args[0]))
	if err != nil{
		return shim.Error("Error GetState house data with " + args[0] + " error, err = "+err.Error())
	}
	return shim.Success(dataBytes)
}

// createHouse create house info with house json string
func (t *HouseRegistryChaincode) createHouse(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	house := &House{}

	//check house data
	err := json.Unmarshal([]byte(args[1]), house)
	if err != nil{
		return shim.Error("Error Unmarshal house data "+err.Error())
	}

	houseID := getHouseID(args[0])
	houseQuery, err := getHouseInfo(houseID, stub)
	if err != nil{
		return shim.Error("createHouse " + err.Error())
	}
	if houseQuery != nil{
		return shim.Error("createHouse GetState data with " + houseID + " already exists")
	}


	err = stub.PutState(houseID, []byte(args[1]))
	if err != nil{
		return shim.Error("Error PutState house data with " + houseID + " error, err = "+err.Error())
	}
	return shim.Success(nil)
}

// updateHouse update house info with house json string
func (t *HouseRegistryChaincode) updateHouse(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	houseUpdate := &House{}

	err := json.Unmarshal([]byte(args[1]), houseUpdate)
	if err != nil{
		return shim.Error("updateHouse Unmarshal house data error, err = "+err.Error())
	}

	houseID := getHouseID(houseUpdate.ID)
	house, err := getHouseInfo(houseID, stub)
	if err != nil{
		return shim.Error("updateHouse " + err.Error())
	}
	if house == nil{
		return shim.Error("updateHouse query house data with " + houseID + " not exists")
	}

	err = stub.PutState(houseID, []byte(args[1]))
	if err != nil{
		return shim.Error("updateHouse PutState house data error, err = "+err.Error())
	}

	return shim.Success(nil)
}

// changeOwner change house owner
func (t *HouseRegistryChaincode) changeOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	houseID := getHouseID(args[0])
	house, err := getHouseInfo(houseID, stub)
	if err != nil{
		return shim.Error("changeOwner " + err.Error())
	}
	if house == nil{
		return shim.Error("createHouse GetState data with " + houseID + " not exists")
	}

	house.Owner = args[1]

	dataSave, err := json.Marshal(*house)
	if err != nil{
		return shim.Error("Error Marshal house data with " + houseID + " error, err = "+err.Error())
	}
	err = stub.PutState(houseID, dataSave)
	if err != nil{
		return shim.Error("Error PutState house data with " + houseID + " error, err = "+err.Error())
	}

	return shim.Success(nil)
}

func (t *HouseRegistryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createHouse" {
		return t.createHouse(stub, args)
	} else if function == "updateHouse"{
		return t.updateHouse(stub, args)
	} else if function == "changeOwner"{
		return t.changeOwner(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"createHouse\", \"updateHouse\", \"changeOwner\"")
}

func (t *HouseRegistryChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "queryInfo" {
		return t.queryHouse(stub, args)
	} else if function == "queryList" {
		return t.queryHouse(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"queryInfo\", \"queryList\"")
}

func main() {
	err := shim.Start(new(HouseRegistryChaincode))
	if err != nil {
		fmt.Printf("Error starting HouseRegistry chaincode: %s", err)
	}
}


func getHouseInfo(id string, stub shim.ChaincodeStubInterface) (*House, error){
	houseID := getHouseID(id)
	dataBytes, err := stub.GetState(houseID)
	if err != nil{
		return nil, errors.New("GetHouseInfo GetState data with " + houseID + " error, err = "+err.Error())
	}
	if dataBytes == nil{
		return nil, nil
	}
	house := &House{}

	err = json.Unmarshal(dataBytes, house)
	if err != nil{
		return nil, errors.New("GetHouseInfo Unmarshal data with " + houseID + " error, err = "+err.Error())
	}
	return house, nil
}

func getHouseID(id string) string{
	return "HouseBaseInfo-" + id
}