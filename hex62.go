package surl

import (
	"math"
	"strings"
)

// Mapping defines algorithm mapping between decimal and string.
type Mapping interface {
	// Itoa decimal to hex string
	Itoa(int64) string
	// Atoi hex string to decimal
	Atoi(string) int64
}

const hexString = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func createKeyMap() map[int]string {
	bytes := []byte(hexString)
	var result = make(map[int]string, len(bytes))
	var l = len(bytes)
	for i := 0; i < l; i++ {
		result[i] = string(bytes[i : i+1])
	}
	return result
}

func createCharMap(keyMap map[int]string) map[string]int {
	var charMap = make(map[string]int, len(keyMap))
	for k, v := range keyMap {
		charMap[v] = k
	}
	return charMap
}

var keyMap = createKeyMap()
var charMap = createCharMap(keyMap)
var keyLen = len(keyMap)

// hex62 is the default transformer of surl.
type hex62 struct{}

// Itoa encodes int64 number to hex string
func (transformer hex62) Itoa(decimal int64) string {
	if decimal == 0 {
		return keyMap[0]
	}
	n := int64(len(keyMap))
	var hex string
	var reminder int64
	for decimal != 0 {
		reminder = decimal % n
		decimal = decimal / n
		hex = keyMap[int(reminder)] + hex
	}
	return hex
}

// Atoi decodes hex string to int64 number
func (transformer hex62) Atoi(hex string) int64 {
	chars := strings.Split(hex, "")
	var ret, base float64
	for i := 0; i < len(chars); i++ {
		base = float64(charMap[chars[i]])
		ret += base * math.Pow(float64(keyLen), float64(len(chars)-i-1))
	}
	return int64(ret)
}
