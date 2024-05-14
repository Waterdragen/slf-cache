package slf_cache

import (
	"github.com/waterdragen/slf-cache/v2/assert"
	"github.com/waterdragen/slf-cache/v2/util"
	"testing"
)

func TestLoadTable(t *testing.T) {
	table := util.LoadTable()
	assert.Ne(len(table), 0)
	//log.Println(table)
}

func TestLoadCorpus(t *testing.T) {
	corpus := util.LoadCorpus("./corpora/shai/trigrams.json")
	assert.Ne(len(corpus), 0)
	//log.Println(corpus)
}
