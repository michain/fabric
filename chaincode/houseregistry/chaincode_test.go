package houseregistry

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)


func TestExample02_Init(t *testing.T) {
	scc := new(HouseRegistryChaincode)
	stub := shim.NewMockStub("house", scc)
	checkInit(t, stub, nil)
}


func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}