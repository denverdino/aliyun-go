package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/denverdino/aliyungo/util"
)

// A Client represents a client of ECS services
type Client struct {
	AccessKeyId     string //Access Key Id
	AccessKeySecret string //Access Key Secret
	debug           bool
	httpClient      *http.Client
	endpoint        string
	version         string
	OwnerId         string
}

// NewClient creates a new instance of ECS client
func (client *Client) Init(endpoint, version, accessKeyId, accessKeySecret string) {
	client.AccessKeyId = accessKeyId
	client.AccessKeySecret = accessKeySecret + "&"
	client.debug = false
	client.httpClient = &http.Client{}
	client.endpoint = endpoint
	client.version = version
}

//NewClient create a new instance of ECS client support ownerID
func (client *Client) InitWithOwnerId(endpoint, version, accessKeyId, accessKeySecret string, ownerId string) {
	client.Init(endpoint, version, accessKeyId, accessKeySecret)
	client.OwnerId = ownerId
}

//set ownerId
func (client *Client) SetOwnerId(ownerId string) {
	client.OwnerId = ownerId
}

// SetEndpoint sets custom endpoint
func (client *Client) SetEndpoint(endpoint string) {
	client.endpoint = endpoint
}

// SetEndpoint sets custom version
func (client *Client) SetVersion(version string) {
	client.version = version
}

// SetAccessKeyId sets new AccessKeyId
func (client *Client) SetAccessKeyId(id string) {
	client.AccessKeyId = id
}

// SetAccessKeySecret sets new AccessKeySecret
func (client *Client) SetAccessKeySecret(secret string) {
	client.AccessKeySecret = secret + "&"
}

// SetDebug sets debug mode to log the request/response message
func (client *Client) SetDebug(debug bool) {
	client.debug = debug
}

// Invoke sends the raw HTTP request for ECS services
func (client *Client) Invoke(action string, args interface{}, response interface{}) error {
	return client.InvokeWithMethod(ECSRequestMethod, action, args, response)
}

/**
 	some interface can only invoke with post request
**/
func (client *Client) InvokeWithMethod(method string, action string, args interface{}, response interface{}) error {

	request := Request{}
	request.initWithOwnerId(client.version, action, client.AccessKeyId, client.OwnerId)

	query := util.ConvertToQueryValues(request)
	util.SetQueryValues(args, &query)

	// Sign request
	signature := util.CreateSignatureForRequest(method, &query, client.AccessKeySecret)

	// Generate the request URL
	requestURL := client.endpoint + "?" + query.Encode() + "&Signature=" + url.QueryEscape(signature)

	var data []byte
	var err error
	if method == POSTRequestMethod && args != nil {
		data, err = json.Marshal(args)
		if err != nil {
			log.Printf("Failed to marshal data")
			return err
		}
	}

	httpReq, err := http.NewRequest(method, requestURL, bytes.NewBuffer(data))

	// TODO move to util and add build val flag
	httpReq.Header.Set("X-SDK-Client", `AliyunGO/`+Version)

	if err != nil {
		return GetClientError(err)
	}

	t0 := time.Now()

	httpResp, err := client.httpClient.Do(httpReq)
	t1 := time.Now()
	if err != nil {
		return GetClientError(err)
	}
	statusCode := httpResp.StatusCode

	if client.debug {
		log.Printf("Invoke %s %s %d (%v)", method, requestURL, statusCode, t1.Sub(t0))
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return GetClientError(err)
	}

	if client.debug {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "    ")
		log.Println(string(prettyJSON.Bytes()))
	}

	if statusCode >= 400 && statusCode <= 599 {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(body, &errorResponse)
		ecsError := &Error{
			ErrorResponse: errorResponse,
			StatusCode:    statusCode,
		}
		return ecsError
	}

	err = json.Unmarshal(body, response)
	//log.Printf("%++v", response)
	if err != nil {
		return GetClientError(err)
	}

	return nil
}

// GenerateClientToken generates the Client Token with random string
func (client *Client) GenerateClientToken() string {
	return util.CreateRandomString()
}

func GetClientErrorFromString(str string) error {
	return &Error{
		ErrorResponse: ErrorResponse{
			Code:    "AliyunGoClientFailure",
			Message: str,
		},
		StatusCode: -1,
	}
}

func GetClientError(err error) error {
	return GetClientErrorFromString(err.Error())
}
