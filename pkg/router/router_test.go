package router

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestTourismService(t *testing.T) {
	t.Run("Service should connect to tourism service", func(t *testing.T) {

		defer gock.Off()

		gock.New("https://tourism.opendatahub.com").
			Get("/v1/Accommodation/2657B7CBCB85380B253D2FBE28AF100E_REDUCED").
			Reply(200).
			JSON(`{"Id": "2657B7CBCB85380B253D2FBE28AF100E_REDUCED"}`)

		path := "/v1/Accommodation/2657B7CBCB85380B253D2FBE28AF100E_REDUCED"
		method := "GET"
		var router = NewDefaultRouter()
		response := router.EntryPoint(path, method)

		assert.Equal(t, 200, response.StatusCode, "Wrong Status Code")
		assert.Contains(t, response.Body, `{"Id": "2657B7CBCB85380B253D2FBE28AF100E_REDUCED"}`)
	})

	t.Run("Service should connect to mobility service", func(t *testing.T) {

		defer gock.Off()
		gock.Observe(gock.DumpRequest)

		gock.New("https://mobility.api.opendatahub.com").
			Get("/v2/tree,node").
			Reply(200).
			JSON(`{"id": "Bicycle"}`)

		path := "/v2/tree,node"
		method := "GET"
		var router = NewDefaultRouter()
		response := router.EntryPoint(path, method)

		assert.Equal(t, 200, response.StatusCode, "Wrong Status Code")
		assert.Contains(t, response.Body, `{"id": "Bicycle"}`)
	})
}
