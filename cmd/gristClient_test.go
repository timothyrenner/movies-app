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
							"IMDB_ID": "tt0084777",
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
							"First_Time": false,
							"uuid": null
						}
					},
					{
						"id": 376,
						"fields": {
							"Name": "Slaughterhouse",
							"IMDB_Link": "https://www.imdb.com/title/tt0093990/",
							"IMDB_ID": "tt0093990",
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
							"First_Time": false,
							"uuid": null
						}
					},
					{
						"id": 372,
						"fields": {
							"Name": "John Wick Chapter 3 - Parabellum",
							"IMDB_Link": "https://www.imdb.com/title/tt6146586/",
							"IMDB_ID": "tt6146586",
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
							"First_Time": false,
							"uuid": null
						}
					}
				]
			}`,
		),
	)

	records, err := client.GetMovieWatchRecords(
		"abc123",
		"Movie_watches",
		nil,
		"-Watched",
		3,
	)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	truth := GristMovieWatchRecords{
		Records: []GristMovieWatchRecord{
			{
				GristRecord: GristRecord{Id: 375},
				Fields: GristMovieWatchFields{
					Name:        "Tenebrae",
					ImdbLink:    "https://www.imdb.com/title/tt0084777/",
					ImdbId:      "tt0084777",
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
					ImdbLink:    "https://www.imdb.com/title/tt0093990/",
					ImdbId:      "tt0093990",
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
					ImdbId:      "tt6146586",
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

func TestUpdateRecords(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder, err := httpmock.NewJsonResponder(200, nil)
	if err != nil {
		t.Errorf("Error creating JSON responder: %v", err)
	}
	httpmock.RegisterResponder(
		"PATCH",
		"https://docs.getgrist.com/api/docs/abc123/tables/Movie_watches/records",
		responder,
	)

	client := NewGristClient("abc-123")
	documentId := "abc123"
	tableId := "Movie_watches"
	records := GristMovieWatchRecords{
		Records: []GristMovieWatchRecord{
			{
				GristRecord: GristRecord{Id: 1},
				Fields:      GristMovieWatchFields{},
			},
		},
	}
	err = client.UpdateMovieWatchRecords(documentId, tableId, &records)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
}

func TestImdbId(t *testing.T) {
	record := GristMovieWatchRecord{
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
	}

	truth := "tt6146586"
	answer := record.ImdbId()

	if truth != answer {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}
