package main

import (
	_"encoding/binary"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
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
	ccID           = "example"
)

// ExampleCC query and transaction arguments
var queryArgs = [][]byte{[]byte("a")}
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


	value := queryCC(client)
	fmt.Println("value1 is", string(value))

	// Move funds
	executeCC(client)


	value = queryCC(client)
	fmt.Println("value2 is", string(value))

	return

	eventID := "test([a-zA-Z]+)"

	// Register chaincode event (pass in channel which receives event details when the event is complete)
	reg, notifier, err := client.RegisterChaincodeEvent(ccID, eventID)
	if err != nil {
		fmt.Println("Failed to register cc event:", err)
	}
	defer client.UnregisterChaincodeEvent(reg)

	// Move funds
	executeCC(client)

	select {
	case ccEvent := <-notifier:
		fmt.Println("Received CC event: %#v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Println(fmt.Sprintf("Did NOT receive CC event for eventId(%s)\n", eventID))
	}
}

func executeCC(client *channel.Client) {
	_, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: txArgs},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Println("Failed to move funds:", err)
	}
}

func queryCC(client *channel.Client) []byte {
	req := channel.Request{ChaincodeID: ccID, Fcn: "query", Args: queryArgs}
	response, err := client.Query(req)

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
