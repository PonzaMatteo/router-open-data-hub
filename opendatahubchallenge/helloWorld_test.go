package main

import "testing"

func TestResponse(t *testing.T) {

	t.Run("successful response", func(t *testing.T) {
		got := getAccommodation()
		want := 200
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}
