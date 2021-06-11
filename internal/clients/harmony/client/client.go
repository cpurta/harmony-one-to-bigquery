package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
)

type latestHeaderResponse struct {
	Result *model.Header `json:"result"`
}

type blockNumberResponse struct {
	Result *model.Block `json:"result"`
}

var _ harmony.HarmonyClient = &harmonyOneClient{}

type harmonyOneClient struct {
	httpClient  *http.Client
	nodeURL     string
	queryID     int
	queryIDLock *sync.Mutex
}

func NewHarmonyOneClient(nodeURL string, httpClient *http.Client) *harmonyOneClient {
	return &harmonyOneClient{
		nodeURL:     nodeURL,
		httpClient:  httpClient,
		queryID:     0,
		queryIDLock: &sync.Mutex{},
	}
}

func (client *harmonyOneClient) GetLatestHeader() (*model.Header, error) {
	var (
		rpcRequest     *http.Request
		rpcResponse    *http.Response
		responseBody   []byte
		headerResponse latestHeaderResponse
		err            error
	)

	if rpcRequest, err = client.buildRequest("hmy_latestHeader", []interface{}{}); err != nil {
		return nil, err
	}

	if rpcResponse, err = client.makeHTTPRequest(rpcRequest); err != nil {
		return nil, err
	}

	defer rpcResponse.Body.Close()

	if responseBody, err = ioutil.ReadAll(rpcResponse.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(responseBody, &headerResponse); err != nil {
		return nil, err
	}

	return headerResponse.Result, nil
}

func (client *harmonyOneClient) GetBlockByNumber(blockNumber int64) (*model.Block, error) {
	var (
		rpcRequest          *http.Request
		rpcResponse         *http.Response
		responseBody        []byte
		blockNumberResponse blockNumberResponse
		err                 error
	)

	if rpcRequest, err = client.buildRequest("hmy_getBlockByNumber", []interface{}{blockNumber, true}); err != nil {
		return nil, err
	}

	if rpcResponse, err = client.makeHTTPRequest(rpcRequest); err != nil {
		return nil, err
	}

	defer rpcResponse.Body.Close()

	if responseBody, err = ioutil.ReadAll(rpcResponse.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(responseBody, &blockNumberResponse); err != nil {
		return nil, err
	}

	return blockNumberResponse.Result, nil
}

func (client *harmonyOneClient) buildRequest(method string, params []interface{}) (*http.Request, error) {
	client.queryIDLock.Lock()

	var (
		requestBody []byte
		requestMap  = map[string]interface{}{
			"jsonrpc": "v1",
			"id":      client.queryID,
			"method":  method,
			"params":  params,
		}
		httpRequest *http.Request
		err         error
	)

	client.queryID += 1

	client.queryIDLock.Unlock()

	if requestBody, err = json.Marshal(requestMap); err != nil {
		return nil, err
	}

	if httpRequest, err = http.NewRequest(http.MethodPost, client.nodeURL, strings.NewReader(string(requestBody))); err != nil {
		return nil, err
	}

	httpRequest.Header.Add("Content-Type", "application/json")

	return httpRequest, nil
}

func (client *harmonyOneClient) makeHTTPRequest(req *http.Request) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)

	if response, err = client.httpClient.Do(req); err != nil {
		return nil, err
	}

	return response, nil
}