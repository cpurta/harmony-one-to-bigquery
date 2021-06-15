package client

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"go.uber.org/zap"
)

// latestHeaderResponse wraps a Header response object into the result object
type latestHeaderResponse struct {
	Result *model.Header `json:"result"`
}

// blockNumberResponse wraps a Block response object into the result object
type blockNumberResponse struct {
	Result *model.Block `json:"result"`
}

var _ harmony.HarmonyClient = &harmonyOneClient{}

// harmonyOneClient is the implementation of the HarmonyOneClient interface.
type harmonyOneClient struct {
	httpClient  *http.Client
	nodeURL     string
	queryID     int
	queryIDLock *sync.Mutex
	logger      *zap.Logger
}

// NewHarmonyOneClient creates a new HarmonyOneClient implementation that will connect
// to a given Harmony One node to pull header and blockchain data.
func NewHarmonyOneClient(nodeURL string, httpClient *http.Client, logger *zap.Logger) *harmonyOneClient {
	return &harmonyOneClient{
		nodeURL:     nodeURL,
		httpClient:  httpClient,
		queryID:     0,
		queryIDLock: &sync.Mutex{},
		logger:      logger,
	}
}

// GetLatestHeader will return the latest block Header that has been submitted to
// the Harmony One blockchain.
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

// getBlockByNumber will return all Block data associated with the given block.
func (client *harmonyOneClient) GetBlockByNumber(blockNumber int64) (*model.Block, error) {
	var (
		rpcRequest          *http.Request
		rpcResponse         *http.Response
		statusCode          int
		sleepDuration       time.Duration
		responseBody        []byte
		blockNumberResponse blockNumberResponse
		err                 error
	)

	if rpcRequest, err = client.buildRequest("hmy_getBlockByNumber", []interface{}{blockNumber, true}); err != nil {
		return nil, err
	}

	for i := 0; statusCode != http.StatusOK; i++ {
		if rpcResponse, err = client.makeHTTPRequest(rpcRequest); err != nil {
			client.logger.Error("recieved error when making getBlockByNumber rpc request", zap.Int64("block_number", blockNumber), zap.Duration("sleep_duration", sleepDuration), zap.Error(err))
		}

		statusCode = rpcResponse.StatusCode

		sleepDuration = time.Second * time.Duration(int64(math.Pow(float64(2), float64(i))))

		time.Sleep(sleepDuration)
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

// buildRequest will build a pointer http.Request based on the Harmony One RPC API
// documentation.
func (client *harmonyOneClient) buildRequest(method string, params []interface{}) (*http.Request, error) {
	client.queryIDLock.Lock()

	var (
		requestBody []byte
		requestMap  = map[string]interface{}{
			"jsonrpc": "2.0",
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

// makeHTTPRequest will send a given http request and return the resulting response
// or an error.
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
