package mapper

type Mapper struct {}

func NewMapper() Mapper {
	return Mapper{}
}

func (m Mapper) Transform(input string) (string, error) {
	return  `
	{
		"id": "2022-05-11 00:00:00.000+0000",
		"start_date": "2022-05-10 00:00:00.000+0000",
		"end_date": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
	}
	`, nil
}
