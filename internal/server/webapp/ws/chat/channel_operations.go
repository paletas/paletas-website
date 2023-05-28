package chat

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	JoinChannelOperationType      = "joinChannel"
	LeaveChannelOperationType     = "leaveChannel"
	ChannelStatusOperationType    = "channelStatus"
	BroadcastMessageOperationType = "broadcastMessage"
)

var operationTypeToStructType = map[string]reflect.Type{
	JoinChannelOperationType:      reflect.TypeOf(JoinChannelOperation{}),
	LeaveChannelOperationType:     reflect.TypeOf(LeaveChannelOperation{}),
	BroadcastMessageOperationType: reflect.TypeOf(BroadcastMessageOperation{}),
	ChannelStatusOperationType:    reflect.TypeOf(ChannelStatusUpdateOperation{}),
}

func UnmarshalOperation(serializedOperation []byte) (IChatOperation, error) {
	var operation BaseOperation
	err := json.Unmarshal(serializedOperation, &operation)
	if err != nil {
		return nil, err
	}

	structType, ok := operationTypeToStructType[operation.Operation]
	if !ok {
		return nil, fmt.Errorf("unknown operation type: %s", operation.Operation)
	}

	operationValue := reflect.New(structType).Interface()
	err = json.Unmarshal(serializedOperation, operationValue)
	if err != nil {
		return nil, err
	}

	return operationValue.(IChatOperation), nil
}

type BaseChannelOperation struct {
	BaseOperation
	Channel string `json:"channel"`
}

func (operation *BaseChannelOperation) GetChannel() string {
	return operation.Channel
}

type BroadcastMessageOperation struct {
	BaseChannelOperation
	Message string `json:"message"`
}

func NewBroadcastMessageOperation(channel string, message string) *BroadcastMessageOperation {
	return &BroadcastMessageOperation{
		BaseChannelOperation: BaseChannelOperation{
			BaseOperation: BaseOperation{
				Operation: BroadcastMessageOperationType,
			},
			Channel: channel,
		},
		Message: message,
	}
}

func (operation *BroadcastMessageOperation) Perform(user *ChatUser) error {
	channel := findChannel(operation.Channel)
	channel.broadcastMessage(operation.Message)
	return nil
}

type ChannelStatusUpdateOperation struct {
	BaseChannelOperation
	UsersOnline int `json:"usersOnline"`
}

func NewChannelStatusUpdateOperation(channel string, usersOnline int) *ChannelStatusUpdateOperation {
	return &ChannelStatusUpdateOperation{
		BaseChannelOperation: BaseChannelOperation{
			BaseOperation: BaseOperation{
				Operation: ChannelStatusOperationType,
			},
			Channel: channel,
		},
		UsersOnline: usersOnline,
	}
}

func (operation *ChannelStatusUpdateOperation) Perform(user *ChatUser) error {
	channel := findChannel(operation.Channel)
	channel.broadcastChannelStatus()
	return nil
}

type JoinChannelOperation struct {
	BaseChannelOperation
}

func NewJoinChannelOperation(channel string) *JoinChannelOperation {
	return &JoinChannelOperation{
		BaseChannelOperation: BaseChannelOperation{
			BaseOperation: BaseOperation{
				Operation: JoinChannelOperationType,
			},
			Channel: channel,
		},
	}
}

func (operation *JoinChannelOperation) Perform(user *ChatUser) error {
	channel := findChannel(operation.Channel)
	room := findChannel(operation.Channel)
	if room != nil {
		room.Join(NewChatUser(user.Username, user.Connection))
	}

	channel.broadcastChannelStatus()

	return nil
}

type LeaveChannelOperation struct {
	BaseChannelOperation
}

func NewLeaveChannelOperation(channel string) *LeaveChannelOperation {
	return &LeaveChannelOperation{
		BaseChannelOperation: BaseChannelOperation{
			BaseOperation: BaseOperation{
				Operation: LeaveChannelOperationType,
			},
			Channel: channel,
		},
	}
}

func (operation *LeaveChannelOperation) Perform(user *ChatUser) error {
	channel := findChannel(operation.Channel)
	room := findChannel(operation.Channel)
	if room != nil {
		room.Leave(NewChatUser(user.Username, user.Connection))
	}

	channel.broadcastChannelStatus()

	return nil
}
