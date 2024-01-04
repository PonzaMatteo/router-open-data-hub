package mapper

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	t.Run("Mapper should map one response from server to configured format", func(t *testing.T) {
		// Arrange:
		var m = NewMapper()

		// configure the mapper ...
		m.AddMapping("evuuid", "id")

		// Act:
		actual, err := m.Transform(`
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

		var m = NewMapper()

		m.AddMapping("evuuid", "id")
		m.AddMapping("evstart", "start_date")
		m.AddMapping("evend", "end_date")

		actual, err := m.Transform(`
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

		var m = NewMapper()

		m.AddMapping("evuuid", "id")
		m.AddMapping("evstart", "start_date")
		m.AddMapping("evend", "end_date")

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
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"			
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Mapping should come from config file", func(t *testing.T) {

		m, err := NewMapperFromFile("../router/config.json", "v2/flat,event")
		assert.NoError(t, err)

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
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e",
			"start_date": "2022-05-10 00:00:00.000+0000",
			"end_date": "2022-05-11 00:00:00.000+0000"			
		}
		`

		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual)
	})

	t.Run("Read input JSON response from file", func(t *testing.T) {

		m, err := NewMapperFromFile("../router/config.json", "v2/flat,event")
		assert.NoError(t, err)

		inputJson, err := readResponseFromFile("response.json")
		assert.NoError(t, err)

		actual, err := m.Transform(inputJson)
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


	t.Run("Mapper should be able to map lists applying the conversion to every element", func(t *testing.T) {
		// Arrange:
		var mapper = NewMapper()

		// how do we represent the intention of mapping every element of the array data?
		 mapper.AddMapping("evuuid", "id") 
		// b. mapper.AddMappingForEachElement("data", "evuuid", "id")
		// c. mapper.AddMappingForEachElement("data",  NewMapper("evuuid", "id"))

		// Act:
		var actual, err = mapper.Transform(`
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

	t.SkipNow()

	t.Run("Mapper should be able to map fields in an array to to the given format", func(t *testing.T) {
		// Arrange:
		var mapper = NewMapper()

		// how do we represent the intention of mapping every element of the array data?
		 mapper.AddMapping("data.evuuid", "data.id")

		// b. mapper.AddMappingForEachElement("data", "evuuid", "id")
		// c. mapper.AddMappingForEachElement("data",  NewMapper("evuuid", "id"))

		// Act:
		var actual, err = mapper.Transform(`
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


	//to work on
	t.Run("Read complex JSON response from file", func(t *testing.T) {

		m, err := NewMapperFromFile("../router/config.json", "v2/flat,event")
		assert.NoError(t, err)

		inputJson, err := readResponseFromFile("complex-response.json")
		assert.NoError(t, err)

		actual, err := m.Transform(inputJson)
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

	//to work on
	t.Run("Read input JSON response from mobility-events file", func(t *testing.T) {

		m, err := NewMapperFromFile("../router/config.json", "v2")
		assert.NoError(t, err)

		inputJson, err := readResponseFromFile("../response-samples/mobility-events.json")
		assert.NoError(t, err)

		actual, err := m.Transform(inputJson)
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
	var out = in
	out = strings.ReplaceAll(out, "\t", "")
	out = strings.ReplaceAll(out, " ", "")
	out = strings.ReplaceAll(out, "\n", "")
	return out
}
