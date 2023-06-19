package btcd

import (
	"net/rpc/jsonrpc"
)

type JsonRpcApi struct {
	BaseUrl string
}

func NewJsonRpcApi(baseUrl string) *JsonRpcApi {
	return &JsonRpcApi{BaseUrl: baseUrl}
}

func (api *JsonRpcApi) GetBlockTemplate() (*BlockTemplate, error) {
	client, err := jsonrpc.Dial("tcp", api.BaseUrl)
	if err != nil {
		return nil, err
	}

	var result BlockTemplate
	err = client.Call("getblocktemplate", []interface{}{}, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (api *JsonRpcApi) StreamBlockTemplates() (chan *BlockTemplate, error) {
	client, err := jsonrpc.Dial("tcp", api.BaseUrl)
	if err != nil {
		return nil, err
	}

	returnChannel := make(chan *BlockTemplate)

	// loop forever, call getblocktemplate and whenever longpollid changes, send the result to the channel
	var lastBlock BlockTemplate
	for {
		var result BlockTemplate
		err = client.Call("getblocktemplate", []interface{}{}, &result)
		if err != nil {
			return nil, err
		}

		if lastBlock.LongPollId != result.LongPollId {
			returnChannel <- &result
			lastBlock = result
		}
	}
}
