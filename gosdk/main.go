package main

import (
	_"encoding/binary"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel"
	orgName        = "org2"
	orgAdmin       = "Admin"
	ordererOrgName = "bcs06071508"
	orgUser		   = "Admin"
	ccID           = "houseregistry"
)

// ExampleCC query and transaction arguments
var queryArgs = [][]byte{[]byte("HouseInfo-1")}
var txArgs = [][]byte{[]byte("a"), []byte("b"), []byte("10")}

func setupAndRun(configOpt core.ConfigProvider, sdkOpts ...fabsdk.Option) {
	//Init the sdk config
	sdk, err := fabsdk.New(configOpt, sdkOpts...)
	if err != nil {
		fmt.Println("Failed to create new SDK:", err)
	}
	defer sdk.Close()
	// ************ setup complete ************** //

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser), fabsdk.WithOrg(orgName))

	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)

	if err != nil {
		fmt.Println("Failed to create new channel client:", err)
	}


	value := queryHouse(client)
	fmt.Println("value1 is", string(value))


}


func queryHouse(client *channel.Client) []byte {
	req := channel.Request{ChaincodeID: ccID, Fcn: "queryHouse", Args: queryArgs}
	response, err := client.Execute(req)

	if err != nil {
		fmt.Println("Failed to query funds:", err)
	}
	fmt.Println(response)

	return response.Payload
}

func main() {
	configPath := "D:/Go-MyPath/src/github.com/michain/fabric/gosdk/config.yaml"
	//End to End testing
	setupAndRun(config.FromFile(configPath))
}
