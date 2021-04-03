package memviz

import (
	"fmt"
	"reflect"
	"unsafe"
)

const (
	inlinable    = true
	ignoredValue = "***"
)

func (m *mapper) mapStruct(structVal reflect.Value) (nodeID, string) {
	uType := structVal.Type()
	id := m.getNodeID(structVal)
	key := getNodeKey(structVal)
	m.nodeSummaries[key] = escapeString(uType.String())

	structTypeName := uType.Name()
	snode := createNode(id, structTypeName)

	for index := 0; index < uType.NumField(); index++ {
		field := structVal.Field(index)
		if !field.CanAddr() {
			// TODO: when does this happen? Can we work around it?
			continue
		}
		anonymous := uType.Field(index).Anonymous
		fieldName := uType.Field(index).Name
		fieldType := uType.Field(index).Type.Name()

		_ = anonymous
		_ = fieldType

		switch skipField("struct", structTypeName, fieldName) {
		case doNotSkip:
			break
		case ignoreCompletely:
			continue
		case ignoreValue:
			snode.addFieldInlined(getStructRef(index), fieldName, ignoredValue)
			continue
		}

		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		fieldID, summary := m.mapValue(field, id, inlinable)

		// if field was inlined (id == 0) then print summary, else just the name and a link to the actual
		if fieldID == 0 {
			snode.addFieldInlined(getStructRef(index), fieldName, interpretValueType(summary, fieldType))
		} else {
			snode.addField(getStructRef(index), fieldName)
			m.addConnection(id, getStructRef(index), fieldID)
		}
	}

	m.addNode(snode)
	return id, m.nodeSummaries[key]
}

func (m *mapper) mapSlice(sliceVal reflect.Value, parentID nodeID, inlineable bool) (nodeID, string) {
	sliceID := m.getNodeID(sliceVal)
	key := getNodeKey(sliceVal)
	sliceType := escapeString(sliceVal.Type().String())
	m.nodeSummaries[key] = sliceType

	if sliceVal.Len() == 0 {
		m.nodeSummaries[key] = sliceType + "\\{\\}"

		if inlineable {
			return 0, m.nodeSummaries[key]
		}

		return m.newBasicNode(sliceVal, m.nodeSummaries[key]), sliceType
	}

	snode := createNode(sliceID, sliceType)

	// sourceID is the nodeID that links will start from
	// if inlined then these come from the parent
	// if not inlined then these come from this node
	sourceID := sliceID
	if inlineable && sliceVal.Len() <= m.inlineableItemLimit {
		//		sourceID = parentID
	}

	length, totalLength := sliceVal.Len(), sliceVal.Len()
	if length > Options().MaxSliceLength {
		length = Options().MaxSliceLength
	}

	for index := 0; index < length; index++ {
		indexID, summary := m.mapValue(sliceVal.Index(index), sliceID, true)
		if indexID != 0 {
			// need pointer to value
			snode.addField(getSliceRef(sliceID, index), fmt.Sprintf("%d", index))
			m.addConnection(sourceID, getSliceRef(sliceID, index), indexID)
		} else {
			// field was inlined so print summary
			snode.addFields(getSliceRef(sliceID, index), fmt.Sprintf("%d", index), getValueRef(sliceID, index), summary)
		}
	}

	if totalLength != length {
		snode.addField(getSliceRef(sliceID, totalLength-1), fmt.Sprintf("%d more ...", (totalLength-length)))
	}

	m.addNode(snode)
	return sliceID, m.nodeSummaries[key]
}

func (m *mapper) mapMap(mapVal reflect.Value, parentID nodeID, inlineable bool) (nodeID, string) {
	// create a string type while escaping graphviz special characters
	mapType := escapeString(mapVal.Type().String())

	nodeKey := getNodeKey(mapVal)

	if mapVal.Len() == 0 {
		m.nodeSummaries[nodeKey] = mapType + "\\{\\}"

		if inlineable {
			return 0, m.nodeSummaries[nodeKey]
		}

		return m.newBasicNode(mapVal, m.nodeSummaries[nodeKey]), mapType
	}

	mapID := m.getNodeID(mapVal)
	var id nodeID
	if inlineable && mapVal.Len() <= m.inlineableItemLimit {
		m.nodeSummaries[nodeKey] = mapType
		id = parentID
	} else {
		id = mapID
	}

	snode := createNode(id, mapType)

	for index, mapKey := range mapVal.MapKeys() {

		if index > Options().MaxMapEntries {
			break
		}

		keyID, keySummary := m.mapValue(mapKey, id, true)
		valueID, valueSummary := m.mapValue(mapVal.MapIndex(mapKey), id, true)

		switch skipField("map", keySummary, "") {
		case doNotSkip:
			break
		case ignoreCompletely:
			continue
		case ignoreValue:
			valueSummary = ignoredValue
		}

		snode.addFields(
			getKeyRef(mapID, index), keySummary,
			getValueRef(mapID, index), valueSummary)

		if keyID != 0 {
			m.addConnection(id, getKeyRef(mapID, index), keyID)
		}
		if valueID != 0 {
			m.addConnection(id, getValueRef(mapID, index), valueID)
		}
	}

	m.addNode(snode)
	return id, m.nodeSummaries[nodeKey]
}

const (
	formatIndex = "%di%d"
	formatKey   = "%dk%d"
	formatValue = "%dv%d"
	portTitle   = "name"
)

func getStructRef(index int) string {
	return fmt.Sprintf("f%d", index)
}

func getSliceRef(sliceID nodeID, index int) string {
	return fmt.Sprintf(formatIndex, sliceID, index)
}

func getKeyRef(sliceID nodeID, index int) string {
	return fmt.Sprintf(formatKey, sliceID, index)
}

func getValueRef(sliceID nodeID, index int) string {
	return fmt.Sprintf(formatValue, sliceID, index)
}
