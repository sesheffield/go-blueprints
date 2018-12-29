package meander_test

import (
	"testing"

	"github.com/sesheffield/go-blueprints/services/meander"
	"github.com/stretchr/testify/assert"
)

func TestCostValues(t *testing.T) {
	assert.Equal(t, 1, int(meander.Cost1))
	assert.Equal(t, 2, int(meander.Cost2))
	assert.Equal(t, 3, int(meander.Cost3))
	assert.Equal(t, 4, int(meander.Cost4))
	assert.Equal(t, 5, int(meander.Cost5))
}

func TestCostString(t *testing.T) {
	assert.Equal(t, "$", meander.Cost1.String())
	assert.Equal(t, "$$", meander.Cost2.String())
	assert.Equal(t, "$$$", meander.Cost3.String())
	assert.Equal(t, "$$$$", meander.Cost4.String())
	assert.Equal(t, "$$$$$", meander.Cost5.String())
}

func TestParseCost(t *testing.T) {
	assert.Equal(t, meander.Cost1, meander.ParseCost("$"))
	assert.Equal(t, meander.Cost2, meander.ParseCost("$$"))
	assert.Equal(t, meander.Cost3, meander.ParseCost("$$$"))
	assert.Equal(t, meander.Cost4, meander.ParseCost("$$$$"))
	assert.Equal(t, meander.Cost5, meander.ParseCost("$$$$$"))
}

func TestParseCostRange(t *testing.T) {
	l := meander.ParseCostRange("$$...$$$")
	assert.Equal(t, meander.Cost2, l.From)
	assert.Equal(t, meander.Cost3, l.To)
	l = meander.ParseCostRange("$...$$$$$")
	assert.Equal(t, meander.Cost1, l.From)
	assert.Equal(t, meander.Cost5, l.To)
}

func TestParseCostRangeString(t *testing.T) {
	assert.Equal(t, "$$...$$$$", (&meander.CostRange{
		From: meander.Cost2,
		To:   meander.Cost4,
	}).String())
}
