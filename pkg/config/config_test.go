package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	t.Run("it should parse yaml file without errors", func(t *testing.T) {
		c, err := FromFile("./testdata/test-config.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, c.Routes)
	})

	t.Run("it should parse json file without errors", func(t *testing.T) {
		c, err := FromFile("./testdata/test-config.json")
		assert.NoError(t, err)
		assert.NotEmpty(t, c.Routes)
	})
}
