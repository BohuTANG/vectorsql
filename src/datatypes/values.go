// Copyright 2019 The OctoSQL Authors.
// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datatypes

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func MakeNull() *Value {
	return &Value{Value: &Value_Null{Null: true}}
}
func ZeroNull() *Value {
	return &Value{Value: &Value_Null{Null: true}}
}

type Phantom struct{}

func MakePhantom() *Value {
	return &Value{Value: &Value_Phantom{Phantom: true}}
}
func ZeroPhantom() *Value {
	return &Value{Value: &Value_Phantom{Phantom: true}}
}

type Int int

func MakeInt(v int) *Value {
	return &Value{Value: &Value_Int{Int: int64(v)}}
}
func ZeroInt() *Value {
	return &Value{Value: &Value_Int{Int: int64(0)}}
}

type Float float64

func MakeFloat(v float64) *Value {
	return &Value{Value: &Value_Float{Float: v}}
}
func ZeroFloat() *Value {
	return &Value{Value: &Value_Float{Float: 0}}
}

func MakeBool(v bool) *Value {
	return &Value{Value: &Value_Bool{Bool: v}}
}
func ZeroBool() *Value {
	return &Value{Value: &Value_Bool{Bool: false}}
}

func MakeString(v string) *Value {
	return &Value{Value: &Value_String_{String_: v}}
}
func ZeroString() *Value {
	return &Value{Value: &Value_String_{String_: ""}}
}

func MakeTime(v time.Time) *Value {
	t, err := ptypes.TimestampProto(v)
	if err != nil {
		panic(err)
	}
	return &Value{Value: &Value_Time{Time: t}}
}
func ZeroTime() *Value {
	return &Value{Value: &Value_Time{Time: &timestamp.Timestamp{}}}
}

func MakeDuration(v time.Duration) *Value {
	return &Value{Value: &Value_Duration{Duration: ptypes.DurationProto(v)}}
}
func ZeroDuration() *Value {
	return &Value{Value: &Value_Duration{Duration: &duration.Duration{}}}
}

func MakeTuple(v ...*Value) *Value {
	tuple := &Tuple{
		Fields: v,
	}
	return &Value{Value: &Value_Tuple{Tuple: tuple}}
}
func ZeroTuple() *Value {
	return &Value{Value: &Value_Tuple{Tuple: &Tuple{
		Fields: nil,
	}}}
}

func MakeObject(v map[string]*Value) *Value {
	object := &Object{
		Fields: v,
	}

	return &Value{Value: &Value_Object{Object: object}}
}
func ZeroObject() *Value {
	return &Value{Value: &Value_Object{Object: &Object{
		Fields: nil,
	}}}
}

// NormalizeType brings various primitive types into the type we want them to be.
// All types coming out of data sources have to be already normalized this way.
func ToValue(value interface{}) *Value {
	switch value := value.(type) {
	case nil:
		return MakeNull()
	case bool:
		return MakeBool(value)
	case int:
		return MakeInt(value)
	case int8:
		return MakeInt(int(value))
	case int32:
		return MakeInt(int(value))
	case int64:
		return MakeInt(int(value))
	case uint8:
		return MakeInt(int(value))
	case uint32:
		return MakeInt(int(value))
	case uint64:
		return MakeInt(int(value))
	case float32:
		return MakeFloat(float64(value))
	case float64:
		return MakeFloat(value)
	case []byte:
		return MakeString(string(value))
	case string:
		return MakeString(value)
	case []interface{}:
		out := make([]*Value, len(value))
		for i := range value {
			out[i] = ToValue(value[i])
		}
		return MakeTuple(out...)
	case map[string]interface{}:
		out := make(map[string]*Value)
		for k, v := range value {
			out[k] = ToValue(v)
		}
		return MakeObject(out)
	case *interface{}:
		if value != nil {
			return ToValue(*value)
		}
		return MakeNull()
	case time.Time:
		return MakeTime(value)
	case time.Duration:
		return MakeDuration(value)
	case struct{}:
		return MakePhantom()
	case *Value:
		return value
	}
	panic(fmt.Sprintf("unreachable:%T", value))
}

func ZeroValue() Value {
	return Value{}
}

func (v Value) AsInt() int {
	return int(v.GetInt())
}

func (v Value) AsFloat() float64 {
	return v.GetFloat()
}

func (v Value) AsBool() bool {
	return v.GetBool()
}

func (v Value) AsString() string {
	return v.GetString_()
}

func (v Value) AsTime() time.Time {
	t, err := ptypes.Timestamp(v.GetTime())
	if err != nil {
		panic(err)
	}
	return t
}

func (v Value) AsDuration() time.Duration {
	d, err := ptypes.Duration(v.GetDuration())
	if err != nil {
		panic(err)
	}
	return d
}

func (v Value) AsSlice() []*Value {
	t := v.GetTuple()
	return t.Fields
}

func (v Value) AsMap() map[string]Value {
	obj := v.GetObject()
	out := make(map[string]Value)
	for k, v := range obj.Fields {
		out[k] = *v
	}
	return out
}

type Type int

const (
	TypeZero Type = iota
	TypeNull
	TypePhantom
	TypeInt
	TypeFloat
	TypeBool
	TypeString
	TypeTime
	TypeDuration
	TypeTuple
	TypeObject
)

// Można na tych Value pod spodem zdefiniowac GetType i użyć wirtualnych metod, a nie type switch
func (v Value) GetType() Type {
	switch v.Value.(type) {
	case *Value_Null:
		return TypeNull
	case *Value_Phantom:
		return TypePhantom
	case *Value_Int:
		return TypeInt
	case *Value_Float:
		return TypeFloat
	case *Value_Bool:
		return TypeBool
	case *Value_String_:
		return TypeString
	case *Value_Time:
		return TypeTime
	case *Value_Duration:
		return TypeDuration
	case *Value_Tuple:
		return TypeTuple
	case *Value_Object:
		return TypeObject
	default:
		return TypeZero
	}
}

func (v Value) Show() string {
	switch v.GetType() {
	case TypeZero:
		return "<zeroValue>"
	case TypeNull:
		return "<null>"
	case TypePhantom:
		return "<phantom>"
	case TypeInt:
		return fmt.Sprint(v.AsInt())
	case TypeFloat:
		return fmt.Sprint(v.AsFloat())
	case TypeBool:
		return fmt.Sprint(v.AsBool())
	case TypeString:
		return fmt.Sprintf("'%s'", v.AsString())
	case TypeTime:
		return v.AsTime().Format(time.RFC3339Nano)
	case TypeDuration:
		return v.AsDuration().String()
	case TypeTuple:
		valueStrings := make([]string, len(v.AsSlice()))
		for i, value := range v.AsSlice() {
			valueStrings[i] = value.Show()
		}
		return fmt.Sprintf("(%s)", strings.Join(valueStrings, ", "))
	case TypeObject:
		pairStrings := make([]string, 0, len(v.AsMap()))
		for k, v := range v.AsMap() {
			pairStrings = append(pairStrings, fmt.Sprintf("%s: %s", k, v.Show()))
		}
		return fmt.Sprintf("{%s}", strings.Join(pairStrings, ", "))
	default:
		panic("invalid type")
	}
}

func (v Value) ToRawValue() interface{} {
	switch v.GetType() {
	case TypeZero:
		return nil
	case TypeNull:
		return nil
	case TypePhantom:
		return struct{}{}
	case TypeInt:
		return v.AsInt()
	case TypeFloat:
		return v.AsFloat()
	case TypeBool:
		return v.AsBool()
	case TypeString:
		return v.AsString()
	case TypeTime:
		return v.AsTime()
	case TypeDuration:
		return v.AsDuration()
	case TypeTuple:
		out := make([]interface{}, len(v.AsSlice()))
		for i, v := range v.AsSlice() {
			out[i] = v.ToRawValue()
		}
		return out
	case TypeObject:
		out := make(map[string]interface{}, len(v.AsMap()))
		for k, v := range v.AsMap() {
			out[k] = v.ToRawValue()
		}
		return out
	default:
		return nil
	}
}
