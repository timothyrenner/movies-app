package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestGetRecords(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewGristClient("abc-123")

	httpmock.RegisterResponder(
		"GET",
		"https://docs.getgrist.com/api/docs/abc123/tables/Movie_watches/records?limit=3&sort=-Watched",
		httpmock.NewStringResponder(
			200,
			`{
				"records": [
					{
						"id": 375,
						"fields": {
							"Name": "Tenebrae",
							"IMDB_Link": "https://www.imdb.com/title/tt0084777/",
							"Watched": 1653609600,
							"Joe_Bob": true,
							"Call_Felissa": false,
							"Beast": false,
							"Godzilla": false,
							"Zombies": false,
							"Slasher": true,
							"Service": [
								"L",
								"Shudder"
							],
							"First_Time": false
						}
					},
					{
						"id": 376,
						"fields": {
							"Name": "Slaughterhouse",
							"IMDB_Link": "https://www.imdb.com/title/tt0093990/?ref_=ext_shr_lnk",
							"Watched": 1653609600,
							"Joe_Bob": true,
							"Call_Felissa": false,
							"Beast": false,
							"Godzilla": false,
							"Zombies": false,
							"Slasher": true,
							"Service": [
								"L",
								"Shudder"
							],
							"First_Time": false
						}
					},
					{
						"id": 372,
						"fields": {
							"Name": "John Wick Chapter 3 - Parabellum",
							"IMDB_Link": "https://www.imdb.com/title/tt6146586/",
							"Watched": 1653436800,
							"Joe_Bob": false,
							"Call_Felissa": true,
							"Beast": false,
							"Godzilla": false,
							"Zombies": false,
							"Slasher": false,
							"Service": [
								"L",
								"Google Play"
							],
							"First_Time": false
						}
					}
				]
			}`,
		),
	)

	records, err := client.GetRecords(
		"abc123",
		"Movie_watches",
		nil,
		"-Watched",
		3,
	)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	truth := GetMovieWatchRecordsResponse{
		Records: []GristMovieWatchRecord{
			{
				GristRecord: GristRecord{Id: 375},
				Fields: GristMovieWatchFields{
					Name:        "Tenebrae",
					ImdbLink:    "https://www.imdb.com/title/tt0084777/",
					FirstTime:   false,
					Watched:     1653609600,
					JoeBob:      true,
					CallFelissa: false,
					Beast:       false,
					Godzilla:    false,
					Zombies:     false,
					Slasher:     true,
					Service:     []string{"L", "Shudder"},
				},
			},
			{
				GristRecord: GristRecord{Id: 376},
				Fields: GristMovieWatchFields{
					Name:        "Slaughterhouse",
					ImdbLink:    "https://www.imdb.com/title/tt0093990/?ref_=ext_shr_lnk",
					FirstTime:   false,
					Watched:     1653609600,
					JoeBob:      true,
					CallFelissa: false,
					Beast:       false,
					Godzilla:    false,
					Zombies:     false,
					Slasher:     true,
					Service:     []string{"L", "Shudder"},
				},
			},
			{
				GristRecord: GristRecord{Id: 372},
				Fields: GristMovieWatchFields{
					Name:        "John Wick Chapter 3 - Parabellum",
					ImdbLink:    "https://www.imdb.com/title/tt6146586/",
					Watched:     1653436800,
					JoeBob:      false,
					CallFelissa: true,
					Beast:       false,
					Godzilla:    false,
					Zombies:     false,
					Slasher:     false,
					FirstTime:   false,
					Service:     []string{"L", "Google Play"},
				},
			},
		},
	}

	if !cmp.Equal(truth, *records) {
		t.Errorf("Expected %v, \n got %v", truth, *records)
		for ii := range truth.Records {
			if !cmp.Equal(truth.Records[ii], records.Records[ii]) {
				t.Errorf(
					"Expected \n %v got \n %v",
					truth.Records[ii],
					records.Records[ii],
				)
			}
		}
	}
}
