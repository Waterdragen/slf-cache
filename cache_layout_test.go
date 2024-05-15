package slf_cache

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/waterdragen/slf-cache/v2/assert"
	"github.com/waterdragen/slf-cache/v2/slf"
	"github.com/waterdragen/slf-cache/v2/util"
)

func TestCacheOneLayoutStats(t *testing.T) {
	layout, err := slf.ReadLayoutFile("./layouts/QWERTY")
	assert.Ok(err)
	corpus := util.LoadCorpus("./corpora/shai/trigrams.json")
	stat := util.AnalyzeTrigrams(&layout, &corpus)
	packedStat := util.PackStats(stat)
	stats := util.NewCachedStat()
	stats.Insert("shai", string(packedStat))

	cachedStats := util.NewCachedStats()
	cachedStats["QWERTY"] = stats
	err = util.WriteCachedStats(cachedStats, "./cached_stats1.json")
	assert.Ok(err)
}

func TestReadOneCachedLayoutStats(t *testing.T) {
	cachedStats, err := util.ReadCachedStats("./cached_stats1.json")
	assert.Ok(err)
	cachedStat, found := cachedStats["QWERTY"]
	assert.Eq(found, true)
	stat, found := cachedStat.Stats["shai"]
	unpackedStat := util.UnpackStats([]byte(stat))
	assert.ApproxEq(unpackedStat[util.Sfb]/2, 0.056465)
	assert.ApproxEq(unpackedStat[util.Alt], 0.18297)
	assert.ApproxEq(unpackedStat[util.InRoll], 0.19032)
	assert.ApproxEq(unpackedStat[util.OutRoll], 0.16129)
}

func TestCacheAllLayoutStats(t *testing.T) {
	corpora := make(map[string][]util.Ngram)
	corporaDir, err := os.ReadDir("./corpora")
	assert.Ok(err)
	for _, corpusEntry := range corporaDir {
		corpusName := corpusEntry.Name()
		corpus := util.LoadCorpus(fmt.Sprintf("./corpora/%s/trigrams.json", corpusName))
		corpora[corpusName] = corpus
	}

	layoutDir, err := os.ReadDir("./layouts")
	assert.Ok(err)

	start := time.Now()
	cachedStats := util.NewCachedStats()
	for _, layoutDirEntry := range layoutDir {
		layoutName := layoutDirEntry.Name()
		layout, err := slf.ReadLayoutFile(fmt.Sprintf("./layouts/%s", layoutName))
		assert.Ok(err)

		// Analyze stat for each corpus
		// Insert stats into CachedStat.Stats
		stats := util.NewCachedStat()
		for corpusName, corpus := range corpora {
			stat := util.AnalyzeTrigrams(&layout, &corpus)
			packedStat := util.PackStats(stat)
			stats.Insert(corpusName, string(packedStat))
		}

		// Insert CachedStat into map[layout]CachedStat
		cachedStats[layoutName] = stats
	}

	err = util.WriteCachedStats(cachedStats, "./cached_stats.json")
	assert.Ok(err)

	assert.Ne(len(cachedStats), 0)
	assert.Eq(len(cachedStats), len(layoutDir))

	elapsed := time.Since(start)
	fmt.Printf("Cache all layout stats took: %s\n", elapsed)
}

func TestReadAllCachedLayoutStats(t *testing.T) {
	cachedStats, err := util.ReadCachedStats("./cached_stats.json")
	assert.Ok(err)

	seen := 0
	for layoutName, cachedStat := range cachedStats {
		stat, found := cachedStat.Stats["shai"]
		assert.Eq(found, true)
		unpackedStat := util.UnpackStats([]byte(stat))
		fmt.Printf("%v: %v\n", layoutName, unpackedStat)

		if seen > 10 {
			break
		}
		seen++
	}
}
