package jsonhelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

type nodeType int

const (
	_ nodeType = iota
	tMap
	tArray
	tNum
	tBool
	tString
)

type Node interface {
	Key(k string) Node
	Index(i int) Node
	IsMap() bool
	IsArray() bool
	IsNum() bool
	IsBool() bool
	IsString() bool
	ToMap() (map[string]Node, error)
	ToArray() ([]Node, error)
	ToInt() (int, error)
	ToInt64() (int64, error)
	ToFloat64() (float64, error)
	ToBool() (bool, error)
	ToString() (string, error)
}

func Parse(data []byte) (Node, error) {
	data = bytes.TrimSpace(data)
	if data[0] == '{' {
		v := map[string]interface{}{}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return createNode(v)
	}
	if data[0] == '[' {
		v := []interface{}{}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return createNode(v)
	}
	return nil, errors.New("Invalid JSON format")
}

func createNode(i interface{}) (Node, error) {
	r := reflect.ValueOf(i)
	switch r.Kind() {
	case reflect.Map:
		n := &nodeMap{
			v: i.(map[string]interface{}),
			c: map[string]Node{},
		}
		return n, nil

	case reflect.Array, reflect.Slice:
		n := &nodeArray{
			v: i.([]interface{}),
			c: map[int]Node{},
		}
		return n, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		n := &nodeValue{
			t: tNum,
			v: i,
		}
		return n, nil

	case reflect.String:
		n := &nodeValue{
			t: tString,
			v: i,
		}
		return n, nil

	case reflect.Bool:
		n := &nodeValue{
			t: tBool,
			v: i,
		}
		return n, nil

	}
	return nil, errors.New("Invalid data type")
}

//nodeError
type nodeError struct {
	e error
}

func (n *nodeError) Key(k string) Node {
	return n
}

func (n *nodeError) Index(i int) Node {
	return n
}

func (n *nodeError) IsMap() bool    { return false }
func (n *nodeError) IsArray() bool  { return false }
func (n *nodeError) IsNum() bool    { return false }
func (n *nodeError) IsBool() bool   { return false }
func (n *nodeError) IsString() bool { return false }

func (n *nodeError) ToMap() (map[string]Node, error) {
	return nil, n.e
}

func (n *nodeError) ToArray() ([]Node, error) {
	return nil, n.e
}

func (n *nodeError) ToInt() (int, error) {
	return 0, n.e
}

func (n *nodeError) ToInt64() (int64, error) {
	return 0, n.e
}

func (n *nodeError) ToFloat64() (float64, error) {
	return 0, n.e
}

func (n *nodeError) ToBool() (bool, error) {
	return false, n.e
}

func (n *nodeError) ToString() (string, error) {
	return "", n.e
}

//nodeMap
type nodeMap struct {
	v map[string]interface{}
	c map[string]Node //cache
}

func (n *nodeMap) Key(k string) Node {
	if nd, ok := n.c[k]; ok {
		return nd
	}
	if v, ok := n.v[k]; ok {
		node, err := createNode(v)
		if err != nil {
			return &nodeError{e: err}
		}
		n.c[k] = node
		return node
	}
	return &nodeError{e: errors.New("Key `" + k + "` not exist")}
}

func (n *nodeMap) Index(i int) Node {
	return &nodeError{e: errors.New("Node is not array")}
}

func (n *nodeMap) IsMap() bool    { return true }
func (n *nodeMap) IsArray() bool  { return false }
func (n *nodeMap) IsNum() bool    { return false }
func (n *nodeMap) IsBool() bool   { return false }
func (n *nodeMap) IsString() bool { return false }

func (n *nodeMap) ToMap() (map[string]Node, error) {
	result := map[string]Node{}
	for k := range n.v {
		result[k] = n.Key(k)
	}
	return result, nil
}

func (n *nodeMap) ToArray() ([]Node, error) {
	return nil, errors.New("Node is not array")
}

func (n *nodeMap) ToInt() (int, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeMap) ToInt64() (int64, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeMap) ToFloat64() (float64, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeMap) ToBool() (bool, error) {
	return false, errors.New("Node is not boolean")
}

func (n *nodeMap) ToString() (string, error) {
	return "", errors.New("Node is not string")
}

//nodeArray
type nodeArray struct {
	v []interface{}
	c map[int]Node //cache
}

func (n *nodeArray) Key(k string) Node {
	return &nodeError{e: errors.New("Node is not map")}
}

func (n *nodeArray) Index(i int) Node {
	if nd, ok := n.c[i]; ok {
		return nd
	}
	if i >= 0 && i < len(n.v) {
		node, err := createNode(n.v[i])
		if err != nil {
			return &nodeError{e: err}
		}
		n.c[i] = node
		return node
	}
	return &nodeError{e: errors.New("Index `" + strconv.Itoa(i) + "` out of range")}
}

func (n *nodeArray) IsMap() bool    { return false }
func (n *nodeArray) IsArray() bool  { return true }
func (n *nodeArray) IsNum() bool    { return false }
func (n *nodeArray) IsBool() bool   { return false }
func (n *nodeArray) IsString() bool { return false }

func (n *nodeArray) ToMap() (map[string]Node, error) {
	return nil, errors.New("Node is not map")
}

func (n *nodeArray) ToArray() ([]Node, error) {
	result := []Node{}
	for i := range n.v {
		result = append(result, n.Index(i))
	}
	return result, nil
}

func (n *nodeArray) ToInt() (int, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeArray) ToInt64() (int64, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeArray) ToFloat64() (float64, error) {
	return 0, errors.New("Node is not number")
}

func (n *nodeArray) ToBool() (bool, error) {
	return false, errors.New("Node is not boolean")
}

func (n *nodeArray) ToString() (string, error) {
	return "", errors.New("Node is not string")
}

//nodeValue
type nodeValue struct {
	t nodeType
	v interface{}
}

func (n *nodeValue) Key(k string) Node {
	return &nodeError{e: errors.New("Node is not map")}
}

func (n *nodeValue) Index(i int) Node {
	return &nodeError{e: errors.New("Node is not array")}
}

func (n *nodeValue) IsMap() bool    { return false }
func (n *nodeValue) IsArray() bool  { return false }
func (n *nodeValue) IsNum() bool    { return n.t == tNum }
func (n *nodeValue) IsBool() bool   { return n.t == tBool }
func (n *nodeValue) IsString() bool { return n.t == tString }

func (n *nodeValue) ToMap() (map[string]Node, error) {
	return nil, errors.New("Node is not map")
}

func (n *nodeValue) ToArray() ([]Node, error) {
	return nil, errors.New("Node is not array")
}

func (n *nodeValue) ToInt() (int, error) {
	if n.IsNum() {
		return int(n.v.(float64)), nil
	}
	return 0, errors.New("Node is not number")
}

func (n *nodeValue) ToInt64() (int64, error) {
	if n.IsNum() {
		return int64(n.v.(float64)), nil
	}
	return 0, errors.New("Node is not number")
}

func (n *nodeValue) ToFloat64() (float64, error) {
	if n.IsNum() {
		return n.v.(float64), nil
	}
	return 0, errors.New("Node is not number")
}

func (n *nodeValue) ToBool() (bool, error) {
	if n.IsBool() {
		return n.v.(bool), nil
	}
	return false, errors.New("Node is not boolean")
}

func (n *nodeValue) ToString() (string, error) {
	if n.IsString() {
		return n.v.(string), nil
	}
	return "", errors.New("Node is not string")
}
