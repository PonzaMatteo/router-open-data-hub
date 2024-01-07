package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndToEndRouter(t *testing.T) {
	t.Run("should contact mobility service and map the response to the configured format", func(t *testing.T) {
		var r, err = NewDefaultRouter()
		assert.NoError(t, err)

		// https://mobility.api.opendatahub.com/v2/flat%2Cevent/%2A/latest?limit=200&offset=0&where=evuuid.eq.53a6343f-e524-51ea-a280-4cc4c1bc7ff3&shownull=false&distinct=true
		response, err := r.RouteRequest("/v2/flat,event/*/latest?limit=200&offset=0&where=evuuid.eq.53a6343f-e524-51ea-a280-4cc4c1bc7ff3&shownull=false&distinct=true", "GET")

		assert.NoError(t, err)
		assert.JSONEq(t, `
		{
			"data": [
				{
					"id": "53a6343f-e524-51ea-a280-4cc4c1bc7ff3",
					"start_date": "2022-05-10 00:00:00.000+0000",
					"end_date": "2022-05-11 00:00:00.000+0000"
				}
			]
		}
		`, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})

	t.Run("should contact tourism service and map the response to the configured format", func(t *testing.T) {
		var r, err = NewDefaultRouter()
		assert.NoError(t, err)

		// https://tourism.opendatahub.com/v1/Event?pagenumber=1&idlist=BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED&removenullvalues=false

		response, err := r.RouteRequest("/v1/Event?pagenumber=1&idlist=BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED&removenullvalues=false", "GET")

		assert.NoError(t, err)
		assert.JSONEq(t, `
		{
			"data": [
				{
					"id": "BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED",
					"start_date": "2022-06-01T00:00:00",
					"end_date":  "2022-06-01T00:00:00"
				}
			]
		}
		`, response.Body)
		assert.Equal(t, 200, response.StatusCode)
	})
}
