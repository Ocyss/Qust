package json

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Data struct {
	data any
}

// New 创建Json类型
func New(body []byte) (*Data, error) {
	j := new(Data)
	err := json.Unmarshal(body, &j.data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Get 根据key获取值
func (j *Data) Get(key string) *Data {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Data{val}
		}
	}
	return &Data{nil}
}

// Gets 根据索引来获取json列表的第index项
func (j *Data) Gets(index int) *Data {
	m, err := j.Maps()
	if err == nil {
		if index >= 0 && index < len(m) {
			return &Data{m[index]}
		}
	}
	return &Data{nil}
}

// Map 判断是不是kv类型
func (j *Data) Map() (map[string]any, error) {
	if v, ok := j.data.(map[string]any); ok {
		return v, nil
	}
	return nil, errors.New("type assertion to map[string]any{} failed")
}

// Maps 判断是不是切片
func (j *Data) Maps() ([]any, error) {
	if v, ok := j.data.([]any); ok {
		return v, nil
	}
	return nil, errors.New("type assertion to []any failed")
}

// Val 获取值
func (j *Data) Val() any {
	return j.data
}

// String 获取字符串
func (j *Data) String() string {
	var out bytes.Buffer
	res, err := json.Marshal(j.data)
	if err != nil {
		return ""
	}
	_ = json.Indent(&out, res, "", "\t")
	return out.String()
}

// Float64 获取浮点数
func (j *Data) Float64() (float64, error) {
	if s, ok := (j.data).(float64); ok {
		return s, nil
	}
	return 0, errors.New("type assertion to float64 failed")
}

// Bool 获取布尔值
func (j *Data) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// Array 获取切片
func (j *Data) Array() ([]any, error) {
	if s, ok := (j.data).([]any); ok {
		return s, nil
	}
	return nil, errors.New("type assertion to []any failed")
}

// Arrays 获取Json列表
func (j *Data) Arrays() ([]*Data, error) {
	var data []*Data
	if s, ok := (j.data).([]any); ok {
		for _, v := range s {
			data = append(data, &Data{v})
		}
		return data, nil
	}
	return nil, errors.New("type assertion to []any failed")
}

// GetVal 直接根据Key获取值
func (j *Data) GetVal(key string) any {
	return j.Get(key).Val()
}

// GetFloat64 直接根据Key获取浮点数
func (j *Data) GetFloat64(key string) float64 {
	v, _ := j.Get(key).Float64()
	return v
}

// GetBool 直接根据Key获取布尔值
func (j *Data) GetBool(key string) bool {
	v, _ := j.Get(key).Bool()
	return v
}

// GetArray  直接根据Key获取切片
func (j *Data) GetArray(key string) []any {
	v, _ := j.Get(key).Array()
	return v
}

// GetArrays 直接根据Key获取Json列表
func (j *Data) GetArrays(key string) []*Data {
	v, _ := j.Get(key).Arrays()
	return v
}
