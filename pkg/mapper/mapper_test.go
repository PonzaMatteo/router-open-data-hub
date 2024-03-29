package mapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	t.Run("Empty mapper should return same output as input", func(t *testing.T) {

		var m = EmptyMapper()
		actual, err := m.MapJSON(`
		{
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapper should map one response from server to configured format", func(t *testing.T) {
		// Arrange:
		var m = EmptyMapper()

		// configure the mapper ...
		m.AddMapping("evuuid", "id")

		// Act:
		actual, err := m.MapJSON(`
		{
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`

		// Assert:
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapper should map response from server to configured format", func(t *testing.T) {

		var m = EmptyMapper()

		m.AddMapping("evuuid", "id")
		m.AddMapping("evstart", "start_date")
		m.AddMapping("evend", "end_date")

		actual, err := m.MapJSON(`
		{
			"evend": "2022-05-11 00:00:00.000+0000",
			"evstart": "2022-05-10 00:00:00.000+0000",
			"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}
		`)

		var expected = `
		{
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"	
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapper should map only required fields", func(t *testing.T) {

		var m = EmptyMapper()

		m.AddMapping("evuuid", "id")
		m.AddMapping("evstart", "start_date")
		m.AddMapping("evend", "end_date")

		actual, err := m.MapJSON(`
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
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"			
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapping should come from config file", func(t *testing.T) {

		m := NewMapper(map[string]string{
			"evuuid":  "id",
			"evstart": "start_date",
			"evend":   "end_date",
		})

		actual, err := m.MapJSON(`
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
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"			
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Read input JSON response from file", func(t *testing.T) {

		m := NewMapper(map[string]string{
			"evuuid":  "id",
			"evstart": "start_date",
			"evend":   "end_date",
		})

		inputJson, err := readResponseFromFile("response.json")
		assert.NoError(t, err)

		actual, err := m.MapJSON(inputJson)
		var expected = `
		{
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"			
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Read complex JSON response", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("data.evuuid", "test.id")
		mapper.AddMapping("data.evstart", "test.start_date")
		mapper.AddMapping("data.evend", "test.end_date")

		var actual, err = mapper.MapJSON(`
		{
			"data": {
				"evend": "2022-05-11 00:00:00.000+0000",
				"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
				"evstart": "2022-05-10 00:00:00.000+0000",
				"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
			}
		  }
		`)

		var expected = `
		{
			"test": 
			{
				"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
				"start_date": "2022-05-10 00:00:00.000+0000",
				"end_date": "2022-05-11 00:00:00.000+0000"			
			}
		  }
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapper should be able to map fields in an array to to the given format", func(t *testing.T) {
		var mapper = EmptyMapper()

		mapper.AddMapping("data.evuuid", "data.id")

		// Act:
		var actual, err = mapper.MapJSON(`
		{
			"data": [
			 { "evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e" },
			 { "evuuid": "74b0c317-2315-4ead-b45f-4acfce220384" }
			]
		}
		`)

		var expected = `
		{
			"data": [
			 { "id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e" },
			 { "id": "74b0c317-2315-4ead-b45f-4acfce220384" }
			]
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapper should be able to map lists applying the conversion to every element", func(t *testing.T) {
		// Arrange:
		var mapper = EmptyMapper()

		mapper.AddMapping("evuuid", "id")

		// Act:
		var actual, err = mapper.MapJSON(`
		[
			{ "evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e" },
			{ "evuuid": "74b0c317-2315-4ead-b45f-4acfce220384" }
		]
		`)

		var expected = `
		[
			{ "id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e" },
			{ "id": "74b0c317-2315-4ead-b45f-4acfce220384" }
		]
		`

		assert.NoError(t, err)
		assertEqualJSON(t, expected, actual)
	})

	t.Run("Read complex JSON response with 2 data", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("data.evuuid", "test.id")
		mapper.AddMapping("data.evstart", "test.start_date")
		mapper.AddMapping("data.evend", "test.end_date")

		var actual, err = mapper.MapJSON(`
		{
			"data": [
				{
					"evend": "2022-05-11 00:00:00.000+0000",
					"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"evstart": "2022-05-10 00:00:00.000+0000",
					"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
				},
				{
					"evend": "2022-05-11 00:00:00.000+0000",
					"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"evstart": "2022-05-10 00:00:00.000+0000",
					"evuuid": "74b0c317-2315-4ead-b45f-4acfce220384"
				}
			]
		  }
		`)

		var expected = `
		{
			"test": [
				{
					"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"start_date": "2022-05-10 00:00:00.000+0000",
					"end_date": "2022-05-11 00:00:00.000+0000"			
				},			
				{
					"id": "74b0c317-2315-4ead-b45f-4acfce220384",
					"start_date": "2022-05-10 00:00:00.000+0000",
					"end_date": "2022-05-11 00:00:00.000+0000"			
				}
			]
		  }
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Read complex JSON response with mutiple nesting", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("data.evuuid", "test.id")
		mapper.AddMapping("data.evstart", "test.start_date")
		mapper.AddMapping("data.evend", "test.end_date")
		mapper.AddMapping("data.evmetadata.placeDe", "test.metadata.german")

		var actual, err = mapper.MapJSON(`
		{
			"offset": 0,
			"data": [
			  {
				"evend": "2022-05-11 00:00:00.000+0000",
				"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
				"evstart": "2022-05-10 00:00:00.000+0000",
				"evtransactiontime": "2022-05-10 18:10:00.663+0000",
				"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
				"evmetadata": {
					"placeDe": "In Richtung Norden zwischen Bozen Süd und Bozen Nord SPERRE ab 21:00 Uhr bis 05:00 Uhr wegen Arbeiten.",
					"placeIt": "In direzione nord tra Bolzano Sud e Bolzano Nord CHIUSURA al traffico dalle ore 21:00 fino alle ore 05:00 a causa di lavori."
				}
			  }
			],
			"limit": 1
		  }
		`)

		var expected = `
		{
			"test": [
				{
					"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"start_date": "2022-05-10 00:00:00.000+0000",
					"end_date": "2022-05-11 00:00:00.000+0000",
					"metadata": {
						"german": "In Richtung Norden zwischen Bozen Süd und Bozen Nord SPERRE ab 21:00 Uhr bis 05:00 Uhr wegen Arbeiten."
					}			
				}
			]
		  }
		`

		assert.NoError(t, err)
		assertEqualJSON(t, expected, actual)
	})

	t.Run("Read input JSON response from mobility-events file", func(t *testing.T) {

		m := NewMapper(map[string]string{
			"data.evuuid":  "data.id",
			"data.evstart": "data.start_date",
			"data.evend":   "data.end_date",
		})

		inputJson, err := readResponseFromFile("../../response-samples/mobility-events.json")
		assert.NoError(t, err)

		actual, err := m.MapJSON(inputJson)
		var expected = `
		{
			"data": [
				{
					"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"start_date": "2022-05-10 00:00:00.000+0000",
					"end_date": "2022-05-11 00:00:00.000+0000"			
				}
			]
		  }
		`

		assert.NoError(t, err)
		assertEqualJSON(t, expected, actual)
	})

	t.Run("Read complex JSON response with error", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("data.evuuid", "data.id")
		mapper.AddMapping("data.evstart", "data.start_date")
		mapper.AddMapping("data.evend", "test.end_date")

		assert.Panics(t, func() {
			mapper.MapJSON(`
			{
				"data": {
					"evend": "2022-05-11 00:00:00.000+0000",
					"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
					"evstart": "2022-05-10 00:00:00.000+0000",
					"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
				}
			  }
			`)
		})

	})

	t.Run("Read complex JSON response that should panic in multiple nested loop", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("another-property.data.evuuid", "another-property.test.id")
		mapper.AddMapping("another-property.data.evstart", "another-property.data.start_date")
		mapper.AddMapping("another-property.data.evend", "another-property.test.end_date")

		assert.Panics(t, func() {
			mapper.MapJSON(`
			{
				"another-property": {
					"data": {
						"evend": "2022-05-11 00:00:00.000+0000",
						"evseriesuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
						"evstart": "2022-05-10 00:00:00.000+0000",
						"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
					}
				}
			  }
			`)
		})

	})

	//NEW

	t.Run("Read tourism data list", func(t *testing.T) {

		var mapper = EmptyMapper()

		mapper.AddMapping("Items.Id", "test.id")
		mapper.AddMapping("Items.DateBegin", "test.start_date")
		mapper.AddMapping("Items.DateEnd", "test.end_date")
		mapper.AddMapping("Items.Active", "test.active")
		mapper.AddMapping("Items.ODHTags", "test.odhtag")
		mapper.AddMapping("Items.Mapping", "test.mapping")

		var actual, err = mapper.MapJSON(`
		{
			"Items": [
			  {
				"Id": "BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED",
				"Active": true,
				"Detail": {
				  "de": {
					"Title": "1000-Stufen-Schlucht im Vinschgau"
				  },
				  "it": {
					"Title": "Venosta: La \"valle dei 1000 gradini\""
				  }
				},
				"DateEnd": "2022-06-01T00:00:00",
				"GpsInfo": [
				  {
					"Gpstype": "position",
					"AltitudeUnitofMeasure": null
				  }
				],
				"LTSTags": null,
				"Mapping": {},
				"ODHTags": [],
				"Latitude": 46.644273,
				"DateBegin": "2022-06-01T00:00:00",
				"Districts": [],
				"EventDate": [
				  {
					"To": "2022-06-01T00:00:00",
					"End": "17:30:00",
					"From": "2022-06-01T00:00:00"
				  }
				]
			  }
			]
		  }
		`)

		var expected = `
		{
			"test": [
				{
				"id": "BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED",
				"end_date": "2022-06-01T00:00:00",
				"start_date": "2022-06-01T00:00:00",
				"active": true,
				"mapping": {},
				"odhtag": []
				}
			]
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})
}

func readResponseFromFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	extension := strings.ToLower(path.Ext(fileName))
	if extension != ".json" {
		return "", fmt.Errorf("unsupported configuration file extension (`%s`): %s", extension, fileName)
	}
	return string(data), nil
}

func assertEqualJSON(t *testing.T, expected string, actual string) {
	t.Helper()

	var expectedString = clean(expected)
	var actualString = clean(actual)
	assert.Equal(t, expectedString, actualString)
}

func clean(in string) string {
	var out any
	err := json.Unmarshal([]byte(in), &out)
	if err != nil {
		panic(err)
	}
	jsonOut, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	return string(jsonOut)
}
