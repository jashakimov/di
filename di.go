package main

import (
	"errors"
	"reflect"
)

type node struct {
	obj interface{}
	key string
}

var container []node

func init() {
	container = make([]node, 0)
}

func Register(obj interface{}) {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		panic("Object must be pointer")
	}
	if !checkObjByKey(obj, "default") {
		panic("Object have already exist")
	}
	container = append(container, node{
		obj: obj,
		key: "default",
	})
}

func Get[T interface{}]() T {
	obj, err := getObjByKey[T]("default")
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func RegisterWithName(obj interface{}, key string) {
	if !checkObjByKey(obj, key) {
		panic("Object have already exist")
	}
	container = append(container, node{
		obj: obj,
		key: key,
	})
}

func GetByName[T interface{}](key string) T {
	obj, err := getObjByKey[T](key)
	if err != nil {
		panic(err.Error())
	}
	return obj
}

func checkObjByKey(obj interface{}, key string) bool {
	for i := range container {
		typesEqual := reflect.TypeOf(obj).String() == reflect.TypeOf(container[i].obj).String()
		prEqual := key == container[i].key
		if typesEqual && prEqual {
			return false
		}
	}
	return true
}

func getObjByKey[T interface{}](key string) (T, error) {
	var empty T

	for i := range container {
		switch container[i].obj.(type) {
		case T:
			if key == container[i].key {
				return container[i].obj.(T), nil
			}
		}
	}
	for i := range container {
		impl := reflect.TypeOf((*T)(nil)).Elem()
		if reflect.TypeOf(container[i].obj).Implements(impl) {
			if key == container[i].key {
				return container[i].obj.(T), nil
			}
		}
	}
	return empty, errors.New("object does not exist in container")
}
