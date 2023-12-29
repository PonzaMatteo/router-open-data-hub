package router

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Run("Router should connect to tourism service", func(t *testing.T) {

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

		// Verify that we don't have pending mocks
		assert.True(t, gock.IsDone())
	})

	t.Run("Router should connect to mobility service", func(t *testing.T) {

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
		assert.True(t, gock.IsDone())
	})

	t.Run("Router should connect to tourism service without mapping", func(t *testing.T) {

		defer gock.Off()
		gock.Observe(gock.DumpRequest)

		gock.New("https://tourism.opendatahub.com").
			Get("/v1/Tag/region").
			Reply(200).
			JSON(`{"Id": "region"}`)

		path := "/v1/Tag/region"
		method := "GET"
		var router = NewDefaultRouter()
		response := router.EntryPoint(path, method)

		assert.Equal(t, 200, response.StatusCode, "Wrong Status Code")
		assert.Contains(t, response.Body, `{"Id": "region"}`)
		assert.True(t, gock.IsDone())
	})

	t.Run("Router should not connect to any service with API that does not exist", func(t *testing.T) {

		defer gock.Off()
		gock.Observe(gock.DumpRequest)

		gock.New("https://tourism.opendatahub.com").
			Get("/v1/Does-Not-Exist").
			Reply(404)

		gock.New("https://mobility.api.opendatahub.com").
			Get("/v1/Does-Not-Exist").
			Reply(404)

		path := "/v1/Does-Not-Exist"
		method := "GET"
		var router = NewDefaultRouter()
		response := router.EntryPoint(path, method)

		assert.Equal(t, 404, response.StatusCode, "Wrong Status Code")
		assert.True(t, gock.IsDone())
	})

	t.Run("Router should not connect to any service with existing keyword but wrong API", func(t *testing.T) {

		defer gock.Off()
		gock.Observe(gock.DumpRequest)

		gock.New("https://mobility.api.opendatahub.com").
			Get("/v2/Does-Not-Exist").
			Reply(404)

		gock.New("https://tourism.opendatahub.com").
			Get("/v2/Does-Not-Exist").
			Reply(404)

		gock.New("https://mobility.api.opendatahub.com").
			Get("/v2/Does-Not-Exist").
			Reply(404)

		path := "/v2/Does-Not-Exist"
		method := "GET"
		var router = NewDefaultRouter()
		response := router.EntryPoint(path, method)

		assert.Equal(t, 404, response.StatusCode, "Wrong Status Code")
		assert.True(t, gock.IsDone())
	})
}
