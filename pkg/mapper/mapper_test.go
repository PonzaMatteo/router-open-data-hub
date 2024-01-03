package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	t.Run("Mapper should map one response from server to configured format", func(t *testing.T) {
		// Arrange:
		var m = NewMapper()

		// configure the mapper ...
		

		// Act:
		actual, err := m.Transform(`
		{
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"id": "2022-05-11 00:00:00.000+0000"
		}
		`

		// Assert:
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)

	} )

	t.Run("Mapper should map response from server to configured format", func(t *testing.T) {
		// Arrange:
		var m = NewMapper()

		// configure the mapper ...
		

		// Act:
		actual, err := m.Transform(`
		{
			"evend": "2022-05-11 00:00:00.000+0000",
			"evstart": "2022-05-10 00:00:00.000+0000",
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"id": "2022-05-11 00:00:00.000+0000",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`

		// Assert:
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)

	} )

	t.Run("Mapper should map only required fields", func(t *testing.T) {
		// Arrange:
		var m = NewMapper()

		// configure the mapper ...
		

		// Act:
		actual, err := m.Transform(`
		{
			"evend": "2022-05-11 00:00:00.000+0000",
			"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"evstart": "2022-05-10 00:00:00.000+0000",
			"evtransactiontime": "2022-05-10 18:10:00.663+0000",
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"id": "2022-05-11 00:00:00.000+0000",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`

		// Assert:
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)

	} )
}
