package util

import (
	"math"
)

const _Interval int = 3

func PackStats(stats []float64) []byte {
	packed := make([]byte, _Interval*len(stats))
	for metricIndex, freq := range stats {
		index := metricIndex * _Interval
		packedFreq := PackFreq(freq)
		copy(packed[index:index+_Interval], packedFreq)
	}
	return packed
}

func UnpackStats(packed []byte) []float64 {
	stats := make([]float64, MetricNum)
	length := len(packed)
	metricIndex := 0
	for index := 0; index < length; index += _Interval {
		packedFreq := packed[index : index+_Interval]
		stats[metricIndex] = UnpackFreq(packedFreq)
		metricIndex++
	}
	return stats
}

func PackBase64(value uint8) byte {
	switch {
	case value < 26:
		return value + 65 // A-Z
	case value < 52:
		return value + 71 // a-z
	case value < 62:
		return value - 4 // 0-9
	case value == 62:
		return '+'
	default:
		return '/'
	}
}

func UnpackBase64(b byte) uint32 {
	ord := uint32(b)
	switch {
	case ord >= 97:
		return ord - 71 // a-z
	case ord >= 65:
		return ord - 65 // A-Z
	case ord >= 48:
		return ord + 4 // 0-9
	case b == '+':
		return 62
	default:
		return 63
	}
}

func PackFreq(f float64) []byte {
	num := uint32(math.Round(f * 100_000.0))
	return []byte{
		PackBase64(uint8(num >> 12 & 0x3f)),
		PackBase64(uint8(num >> 6 & 0x3f)),
		PackBase64(uint8(num & 0x3f)),
	}
}

func UnpackFreq(chars []byte) float64 {
	num := UnpackBase64(chars[0])<<12 | UnpackBase64(chars[1])<<6 | UnpackBase64(chars[2])
	return float64(num) / 100_000.0
}
