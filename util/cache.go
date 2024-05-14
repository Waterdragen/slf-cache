package util

import (
	"encoding/json"
	"os"
)

type CachedStat struct {
	// Reserve a `Sum` field for future uses
	Stats map[string]string `json:"stats"`
}

func NewCachedStat() CachedStat {
	return CachedStat{Stats: make(map[string]string)}
}

func (s CachedStat) Insert(corpusName, stat string) {
	s.Stats[corpusName] = stat
}

func (s CachedStat) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"stats": s.Stats,
	}
	return json.Marshal(data)
}

func NewCachedStats() map[string]CachedStat {
	return make(map[string]CachedStat)
}

func WriteCachedStats(cachedStats map[string]CachedStat, path string) error {
	b, err := json.MarshalIndent(cachedStats, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, os.ModePerm)
}

func ReadCachedStats(path string) (map[string]CachedStat, error) {
	b, err := os.ReadFile(path)
	cachedStats := make(map[string]CachedStat)
	if err != nil {
		return cachedStats, err
	}
	err = json.Unmarshal(b, &cachedStats)
	return cachedStats, err
}
