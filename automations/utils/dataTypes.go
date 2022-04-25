package utils

import (
	"errors"
	"reflect"
	"time"
)

type IncomingMessage struct {
	Payload   map[string]interface{} `json:"payload"`
	UserId    string                 `json:"userId,omitempty"`
	AppId     string                 `json:"appId,omitempty"`
	ArrivedAt time.Time              `json:"arrivedAt"`
}

type Automation struct {
	Id            string         `json:"id" bson:"_id,omitempty"`
	AppId         string         `json:"appId" bson:"app_id,omitempty"`
	UserId        string         `json:"userId" bson:"user_id,omitempty"`
	Description   string         `json:"description" bson:"description"`
	Type          OperationTypes `json:"type" bson:"type,omitempty"`
	FirstOperand  AttrOperand    `json:"firstOperand" bson:"first_operand,omitempty"`
	SecondOperand interface{}    `json:"secondOperand" bson:"second_operand,omitempty"`
	CreatedAt     time.Time      `json:"createdAt" bson:"created_at,omitempty"`
	ModifiedAt    time.Time      `json:"modifiedAt" bson:"modified_at,omitempty"`
}

type AttrOperand map[string]interface{}

//type ConstOperand interface{}

type OperationTypes string

const (
	lt OperationTypes = "lt" //Less than
	gt OperationTypes = "gt" //Greater than
	eq OperationTypes = "eq" //Equal
)

func (a Automation) Check(msg IncomingMessage) (bool, error) {
	firstOperand, err := getAttr(a.FirstOperand, msg.Payload)
	if err != nil {
		return false, err
	}
	if reflect.TypeOf(firstOperand) != reflect.TypeOf(a.SecondOperand) {
		return false, errors.New("operands are from different type")
	}
	switch a.Type {
	//TODO: fix with generics.
	case lt:
		switch firstOperand.(type) {
		case float64:
			return firstOperand.(float64) < a.SecondOperand.(float64), nil
		case int:
			return firstOperand.(int) < a.SecondOperand.(int), nil
		}
	case gt:
		switch firstOperand.(type) {
		case float64:
			return firstOperand.(float64) > a.SecondOperand.(float64), nil
		case int:
			return firstOperand.(int) > a.SecondOperand.(int), nil
		}
	case eq:
		return firstOperand == a.SecondOperand, nil
	}
	return false, errors.New("unknown error")
}

func getAttr(path map[string]interface{}, data map[string]interface{}) (interface{}, error) {
	for key := range path {
		switch path[key].(type) {
		case AttrOperand:
			if _, ok := data[key]; ok {
				return getAttr(path[key].(AttrOperand), data[key].(map[string]interface{}))
			} else {
				return 0, errors.New("invalid operand path inside message")
			}
		case interface{}:
			if _, ok := data[key]; ok {
				return data[key], nil
			} else {
				return 0, errors.New("invalid operand path inside message")
			}
		default:
			return 0, errors.New("unknown error")
		}
	}
	return 0, errors.New("unknown error")
}
