package common

const (
	LEVELDB_TYPE_HASH      = 'h'
	LEVELDB_TYPE_SET       = 's'
	LEVELDB_TYPE_SORTEDSET = 'z'
)

func EncodeHashKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = LEVELDB_TYPE_HASH
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func EncodeSetKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = LEVELDB_TYPE_SET
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func EncodeSortedSetKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = LEVELDB_TYPE_SORTEDSET
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func DecodeHashKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != LEVELDB_TYPE_HASH {
		return nil, nil, false
	}
	nameLen := int(data[1])
	if len(data)-4 < nameLen {
		return nil, nil, false
	}
	name = data[2 : 2+nameLen]
	key = data[3+nameLen:]
	return name, key, true
}

func DecodeSetKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != LEVELDB_TYPE_SET {
		return nil, nil, false
	}
	nameLen := int(data[1])
	if len(data)-4 < nameLen {
		return nil, nil, false
	}
	name = data[2 : 2+nameLen]
	key = data[3+nameLen:]
	return name, key, true
}

func DecodeSortedSetKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != LEVELDB_TYPE_SORTEDSET {
		return nil, nil, false
	}
	nameLen := int(data[1])
	if len(data)-4 < nameLen {
		return nil, nil, false
	}
	name = data[2 : 2+nameLen]
	key = data[3+nameLen:]
	return name, key, true
}
