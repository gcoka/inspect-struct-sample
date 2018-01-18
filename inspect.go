package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// Sample1 simple struct
type Sample1 struct {
	BoolVal bool
	IntVal  int
	StrVal  string
}

// Sample2 pointer fields struct
type Sample2 struct {
	BoolRef *bool
	IntRef  *int
	StrRef  *string
}

func main() {

	fmt.Println("inspect for Sample1 zero value")
	s1zero := Sample1{}
	inspect(s1zero)
	fmt.Println("inspect for Sample2 zero value")
	s2zero := Sample2{}
	inspect(s2zero)
}

func inspect(s interface{}) {

	// declair sub-func for recursive call
	var inspectSub func(interface{}, int)
	inspectSub = func(s interface{}, level int) {
		val := reflect.ValueOf(s)

		printValue(val, level)
		if val.Kind() == reflect.Struct {
			t := reflect.TypeOf(s)
			for i := 0; i < t.NumField(); i++ {
				fmt.Println(i)

				// sf := t.Field(i)
				sfval := val.Field(i)
				if sfval.CanInterface() {
					inspectSub(sfval.Interface(), level+1)
				}
			}
		}
	}

	inspectSub(s, 0)
}

func printValue(val reflect.Value, level int) {
	printWithIndent(fmt.Sprintf("Kind: %v, CanSet: %v", val.Kind(), val.CanSet()), level)
	if val.CanInterface() {
		i := val.Interface()
		printWithIndent(fmt.Sprintf("InnerValue: %v", i), level)
		printWithIndent(fmt.Sprintf("Type: %v", reflect.TypeOf(i)), level)
	}
}

const indentSize = 4

func printWithIndent(s string, level int) {
	spaces := make([]string, 0)
	for i := 0; i < level*indentSize; i++ {
		spaces = append(spaces, " ")
	}
	fmt.Print(strings.Join(spaces, ""))
	fmt.Print(s + "\n")
}

// setVal はreflect.Valueにvalを代入する
func setVal(target reflect.Value, val interface{}) {
	if !target.CanSet() {
		fmt.Printf("target{%v}.CanSet() == false, Kind() = %v \n", target, target.Kind())
		return
	}

	// 対象のStructモデルの型で分岐
	i := target.Interface()
	switch i.(type) {
	case string:
		target.SetString(val.(string))
	case int64, int, int32, int16, int8:
		switch val.(type) {
		case int64:
			target.SetInt(val.(int64))
		case float64:
			target.SetInt(int64(val.(float64)))
		}
	case uint64, uint, uint32, uint16, uint8:
		switch val.(type) {
		case int64:
			target.SetUint(uint64(val.(int64)))
		case float64:
			target.SetUint(uint64(val.(float64)))
		}
	case bool:
		target.SetBool(val.(bool))
	case *bool:
		// fmt.Printf("target: %v *bool!! setting val: %v \n", i, val)
		// elm := reflect.ValueOf(&i).Elem()
		// elm.Set(reflect.ValueOf(val))
		// x := &target
		// *x = elm
		// fmt.Printf("target is set: %v type: %v\n", i, reflect.TypeOf(i))
		// fmt.Printf("elm: %v, target %v\n", elm, target)
		fmt.Printf("Value Kind: %v\n", target.Kind())
		fmt.Printf("Value CanSet: %v\n", target.CanSet())
		fmt.Printf("Value CanAddr: %v\n", target.CanAddr())
		// z := reflect.New(reflect.TypeOf(target.Type))
		target.Set(reflect.ValueOf(val))
	case time.Time:
		if tStr, ok := val.(string); ok {
			t, err := time.Parse("2006-01-02 15:04:05", tStr)
			if err != nil {
				log.Printf("'%v' cannot assign to time.Time", val)
				return
			}
			target.Set(reflect.ValueOf(t))
		}

	default:
		fmt.Printf("setVal doesn't support %v ValueKind: %v\n", target, target.Kind())
	}
}
