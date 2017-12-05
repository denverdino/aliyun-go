package ecs

import "github.com/denverdino/aliyungo/common"

//Modify with your Access Key Id and Access Key Secret

const (
	TestAccessKeyId     = "MY_ACCESS_KEY_ID"
	TestAccessKeySecret = "MY_ACCESS_KEY_SECRET"
	TestInstanceId      = "MY_TEST_INSTANCEID"
	TestSecurityGroupId = "MY_TEST_SECURITY_GROUP_ID"
	TestImageId         = "MY_IMAGE_ID"
	TestAccountId       = "MY_TEST_ACCOUNT_ID" //Get from https://account.console.aliyun.com
	TestRegionID        = common.APNorthEast1
	TestInstanceType    = "ecs.n4.large"
	TestVSwitchID       = "MY_TEST_VSWITCHID"

	TestIAmRich = false
	TestQuick   = false
)

var testClient *Client

func NewTestClient() *Client {
	if testClient == nil {
		testClient = NewClient(TestAccessKeyId, TestAccessKeySecret)
	}
	return testClient
}

var testDebugClient *Client

func NewTestClientForDebug() *Client {
	if testDebugClient == nil {
		testDebugClient = NewClient(TestAccessKeyId, TestAccessKeySecret)
		testDebugClient.SetDebugMode(true)
	}
	return testDebugClient
}

var testLocationClient *Client

func NetTestLocationClientForDebug() *Client {
	if testLocationClient == nil {
		testLocationClient = NewECSClient(TestAccessKeyId, TestAccessKeySecret, TestRegionID)
		testLocationClient.SetDebugMode(true)
	}

	return testLocationClient
}
