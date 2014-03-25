package mustache

import (
	"reflect"
)

func upscope(d interface{}, i interface{}) map[string]interface{} {
	dot, ok := d.(map[string]interface{})
	if !ok {
		dot = make(map[string]interface{})
	} else {
		if dot == nil {
			dot = make(map[string]interface{})
		}
	}
	var changes *mustacheRecord
	subject := i

	switch {
	case mapType(subject):
		changes = upscopeMap(dot, subject)
	case structType(subject):
		changes = upscopeStruct(dot, subject)
	default:
		dot["mustacheItem"] = i
		return dot
	}
	if list, ok := dot["mustacheScopeList"]; ok {
		if records, ok := list.([]*mustacheRecord); ok {
			dot["mustacheScopeList"] = append(records, changes)
		} else {
			dot["mustacheScopeList"] = []*mustacheRecord{changes}
		}
	} else {
		dot["mustacheScopeList"] = []*mustacheRecord{changes}
	}

	return dot
}

func mapType(i interface{}) bool {
	if i != nil {
		return reflect.TypeOf(i).Kind() == reflect.Map
	}
	return false
}

func structType(i interface{}) bool {
	it := reflect.TypeOf(i)
	if it.Kind() == reflect.Ptr {
		it = it.Elem()
	}
	return it.Kind() == reflect.Struct
}
func upscopeMap(dot map[string]interface{}, subject interface{}) *mustacheRecord {
	changes := &mustacheRecord{
		replaced: map[string]interface{}{},
	}

	subjectValue := reflect.ValueOf(subject)
	if subjectValue.Type().Kind() != reflect.Map {
		return nil
	}
	for _, key := range subjectValue.MapKeys() {
		keyName := key.String()
		if val, ok := dot[keyName]; ok {
			changes.replaced[keyName] = val
		} else {
			changes.added = append(changes.added, keyName)
		}
		dot[key.String()] = subjectValue.MapIndex(key).Interface()
	}
	return changes
}

func upscopeStruct(dot map[string]interface{}, subject interface{}) *mustacheRecord {
	return nil
}

func downscope(dot map[string]interface{}) map[string]interface{} {
	record, ok := dot["mustacheScopeList"]
	if !ok {
		return dot
	}
	records, ok := record.([]*mustacheRecord)
	if !ok {
		return dot
	}
	if len(records) > 0 {
		subject := records[len(records)-1]

		if subject == nil {
			return dot
		}
		for _, added := range subject.added {
			delete(dot, added)
		}

		for key, previous := range subject.replaced {
			dot[key] = previous
		}
	}

	return dot
}

type mustacheRecord struct {
	replaced map[string]interface{}
	added    []string
}
