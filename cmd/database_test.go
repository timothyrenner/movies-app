package cmd

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/database"
)

func setupDatabase() (*sql.DB, *migrate.Migrate) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("Encountered error opening in-memory database: %v", err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Panicf(
			"Encountered error creating driver for in-memory database: %v", err,
		)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations", "sqlite3", driver,
	)
	if err != nil {
		log.Panicf("Encountered error creating migration: %v", err)
	}
	if err = m.Up(); err != nil {
		log.Panicf("Encountered error running migration: %v", err)
	}

	return db, m
}

func teardownDatabase(db *sql.DB, m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Panicf("Encountered error tearing down database: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Panicf("Encountered error closing database: %v", err)
	}
}

func omdbSampleMovie() *OmdbMovieResponse {
	return &OmdbMovieResponse{
		Title:    "Tenebrae",
		Year:     "1982",
		Rated:    "R",
		Released: "17 Feb 1984",
		Runtime:  "101 min",
		Genre:    "Horror, Mystery, Thriller",
		Director: "Dario Argento",
		Writer:   "Dario Argento",
		Actors:   "Anthony Franciosa, Giuliano Gemma, John Saxon",
		Plot:     "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.",
		Language: "Italian, Spanish",
		Country:  "Italy",
		Awards:   "N/A",
		Poster:   "https://m.media-amazon.com/images/M/MV5BOTRmNGQ5NTAtNGEzYS00Mjk5LThiZDQtOTk4YTEzNTE1MGZkXkEyXkFqcGdeQXVyNjc1NTYyMjg@._V1_SX300.jpg",
		Ratings: []Rating{
			{Source: "Internet Movie Database", Value: "7.0/10"},
			{Source: "Rotten Tomatoes", Value: "77%"},
			{Source: "Metacritic", Value: "83/100"},
		},
		Metascore:  "83",
		ImdbRating: "7.0",
		ImdbVotes:  "23,156",
		ImdbID:     "tt0084777",
		Type:       "movie",
		DVD:        "20 Sep 2016",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
}

func sampleMovieWatchPage() *MovieWatchPage {
	return &MovieWatchPage{
		Title:       "Tenebrae",
		FileTitle:   "Tenebrae",
		Watched:     "2022-05-27",
		ImdbLink:    "https://www.imdb.com/title/tt0084777/",
		ImdbId:      "tt0084777",
		FirstTime:   false,
		JoeBob:      true,
		CallFelissa: false,
		Beast:       false,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     true,
		WallpaperFu: false,
		Service:     "Shudder",
		Notes:       "",
	}
}

func sampleMoviePage() *MoviePage {
	return &MoviePage{
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		Genres:         []string{"Horror", "Mystery", "Thriller"},
		Directors:      []string{"Dario Argento"},
		Actors:         []string{"Anthony Franciosa", "Giuliano Gemma", "John Saxon"},
		Writers:        []string{"Dario Argento"},
		Year:           1982,
		Rating:         "R",
		Released:       "1984-02-17",
		RuntimeMinutes: 101,
		Plot:           "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.",
		Country:        "Italy",
		Language:       "Italian, Spanish",
		BoxOffice:      "",
		Production:     "",
		CallFelissa:    false,
		Slasher:        true,
		Zombies:        false,
		Beast:          false,
		Godzilla:       false,
		WallpaperFu:    false,
	}
}

func sampleReviewPage() *MovieReviewPage {
	return &MovieReviewPage{
		MovieTitle: "Things",
		Review:     "YOU HAVE JUST EXPERIENCED ... THINGS",
		Liked:      false,
	}
}

func TestCreateInsertMovieParams(t *testing.T) {
	moviePage := sampleMoviePage()

	insertMovieParams, err := CreateInsertMovieParams(moviePage)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	truth := database.InsertMovieParams{
		Uuid:           insertMovieParams.Uuid,
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		ImdbID:         "tt0084777",
		Year:           1982,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes: sql.NullInt64{Int64: 101, Valid: true},
		Plot:           sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:        sql.NullString{String: "Italy", Valid: true},
		Language:       sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    0,
		Beast:          0,
		Slasher:        1,
		Godzilla:       0,
		WallpaperFu:    sql.NullBool{Bool: false, Valid: true},
	}

	if !cmp.Equal(truth, *insertMovieParams) {
		t.Errorf("Expected %v \n got %v", truth, *insertMovieParams)
	}

	// Now test when there's a null runtime.
	prey := MoviePage{
		Title:          "Prey",
		ImdbLink:       "https://www.imdb.com/title/tt11866324/",
		Genres:         []string{"Action", "Drama", "Horror"},
		Directors:      []string{"Dan Trachtenberg"},
		Actors:         []string{"Amber Midthunder", "Dane DiLiegro", "Harlan Blayne Kytwayhat"},
		Writers:        []string{"Patrick Aison"},
		Year:           2022,
		Rating:         "R",
		Released:       "2022-08-05",
		RuntimeMinutes: 0,
		Plot:           "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.",
		Country:        "United States",
		Language:       "English",
		BoxOffice:      "",
		Production:     "",
		CallFelissa:    false,
		Slasher:        false,
		Zombies:        false,
		Beast:          true,
		Godzilla:       false,
		WallpaperFu:    false,
	}

	preyRow, err := CreateInsertMovieParams(&prey)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	preyTruth := database.InsertMovieParams{
		Uuid:           preyRow.Uuid,
		Title:          "Prey",
		ImdbLink:       "https://www.imdb.com/title/tt11866324/",
		ImdbID:         "tt11866324",
		Year:           2022,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "2022-08-05", Valid: true},
		RuntimeMinutes: sql.NullInt64{Int64: 0, Valid: false},
		Plot:           sql.NullString{String: "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.", Valid: true},
		Country:        sql.NullString{String: "United States", Valid: true},
		Language:       sql.NullString{String: "English", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    0,
		Beast:          1,
		Slasher:        0,
		Godzilla:       0,
		WallpaperFu:    sql.NullBool{Bool: false, Valid: true},
	}

	if !cmp.Equal(preyTruth, *preyRow) {
		t.Errorf("Expected \n%v, got \n%v", preyTruth, *preyRow)
	}
}

func TestCreateInsertMovieWatchParams(t *testing.T) {
	movieWatchPage := sampleMovieWatchPage()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieWatchParams(movieWatchPage, movieUuid)
	// Test with null notes.
	truth := database.InsertMovieWatchParams{
		Uuid:       answer.Uuid,
		MovieUuid:  sql.NullString{String: movieUuid, Valid: true},
		MovieTitle: sql.NullString{String: "Tenebrae", Valid: true},
		ImdbID:     "tt0084777",
		Watched:    sql.NullString{String: "2022-05-27", Valid: true},
		Service:    "Shudder",
		FirstTime:  0,
		JoeBob:     1,
		Notes:      sql.NullString{String: "", Valid: false},
	}

	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}

	movieWatchPageWithNotes := sampleMovieWatchPage()
	movieWatchPageWithNotes.Notes = "Great flick"
	answerNotes := CreateInsertMovieWatchParams(movieWatchPageWithNotes, movieUuid)
	truthNotes := database.InsertMovieWatchParams{
		Uuid:       answerNotes.Uuid,
		MovieUuid:  sql.NullString{String: movieUuid, Valid: true},
		MovieTitle: sql.NullString{String: "Tenebrae", Valid: true},
		ImdbID:     "tt0084777",
		Watched:    sql.NullString{String: "2022-05-27", Valid: true},
		Service:    "Shudder",
		FirstTime:  0,
		JoeBob:     1,
		Notes:      sql.NullString{String: "Great flick", Valid: true},
	}
	if !cmp.Equal(truthNotes, *answerNotes) {
		t.Errorf("Expected \n%v, got \n%v", truthNotes, *answerNotes)
	}
}

func TestCreateInsertMovieGenreParams(t *testing.T) {
	moviePage := sampleMoviePage()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieGenreParams(moviePage, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []database.InsertMovieGenreParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Horror",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Mystery",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Thriller",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestTextToNullString(t *testing.T) {
	na := "N/A"
	naTruth := sql.NullString{}
	naAnswer := textToNullString(na)
	if !cmp.Equal(naTruth, naAnswer) {
		t.Errorf("Expected %v, got %v", naTruth, naAnswer)
	}

	empty := ""
	emptyTruth := sql.NullString{}
	emptyAnswer := textToNullString(empty)
	if !cmp.Equal(emptyTruth, emptyAnswer) {
		t.Errorf("Expected %v, got %v", emptyTruth, emptyAnswer)
	}

	notEmpty := "R"
	notEmptyTruth := sql.NullString{String: "R", Valid: true}
	notEmptyAnswer := textToNullString(notEmpty)
	if !cmp.Equal(notEmptyTruth, notEmptyAnswer) {
		t.Errorf("Expected %v, got %v", notEmptyTruth, notEmptyAnswer)
	}
}

func TestCreateInsertMovieActorParams(t *testing.T) {
	moviePage := sampleMoviePage()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieActorParams(moviePage, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []database.InsertMovieActorParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Anthony Franciosa",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Giuliano Gemma",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "John Saxon",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieDirectorParams(t *testing.T) {
	moviePage := sampleMoviePage()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieDirectorParams(moviePage, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []database.InsertMovieDirectorParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieWriterParams(t *testing.T) {
	moviePage := sampleMoviePage()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieWriterParams(moviePage, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []database.InsertMovieWriterParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieRatingParams(t *testing.T) {
	movieRatings := []Rating{
		{
			Source: "Internet Movie Database",
			Value:  "7.0/10.0",
		}, {
			Source: "Rotten Tomatoes",
			Value:  "77%",
		}, {
			Source: "Metacritic",
			Value:  "83/100",
		},
	}
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieRatingParams(movieRatings, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}

	truth := []database.InsertMovieRatingParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Internet Movie Database",
			Value:     "7.0/10.0",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Rotten Tomatoes",
			Value:     "77%",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Metacritic",
			Value:     "83/100",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, answer)
	}
}

func TestCreateInsertMovieReviewParams(t *testing.T) {
	movieReview := sampleReviewPage()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieReviewParams(movieReview, movieUuid)
	truth := &database.InsertReviewParams{
		Uuid:       answer.Uuid,
		MovieUuid:  movieUuid,
		MovieTitle: "Things",
		Review:     "YOU HAVE JUST EXPERIENCED ... THINGS",
		Liked:      0,
	}

	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, answer)
	}
}

func TestInsertMovieDetails(t *testing.T) {
	db, m := setupDatabase()
	defer teardownDatabase(db, m)

	moviePage := sampleMoviePage()
	movieRatings := []Rating{
		{
			Source: "Internet Movie Database",
			Value:  "7.0/10.0",
		}, {
			Source: "Rotten Tomatoes",
			Value:  "77%",
		}, {
			Source: "Metacritic",
			Value:  "83/100",
		},
	}

	queries := database.New(db)
	ctx := context.Background()
	answer, err := InsertMovieDetails(
		db,
		ctx,
		queries,
		moviePage,
		movieRatings,
	)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	movieRowAnswer, err := queries.GetMovie(ctx, answer.Movie)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	movieRowTruth := database.Movie{
		Uuid:            answer.Movie,
		CreatedDatetime: movieRowAnswer.CreatedDatetime,
		Title:           "Tenebrae",
		ImdbLink:        "https://www.imdb.com/title/tt0084777/",
		ImdbID:          "tt0084777",
		Year:            1982,
		Rated:           sql.NullString{String: "R", Valid: true},
		Released:        sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes:  sql.NullInt64{Int64: 101, Valid: true},
		Plot:            sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:         sql.NullString{String: "Italy", Valid: true},
		Language:        sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:       sql.NullString{String: "", Valid: false},
		Production:      sql.NullString{String: "", Valid: false},
		CallFelissa:     0,
		Slasher:         1,
		Zombies:         0,
		Beast:           0,
		Godzilla:        0,
		WallpaperFu:     sql.NullBool{Bool: false, Valid: true},
	}

	if !cmp.Equal(movieRowTruth, movieRowAnswer) {
		t.Errorf("Expected \n%v, got \n%v", movieRowTruth, movieRowAnswer)
	}

	// Movie genre rows.
	if len(answer.Genre) != 3 {
		t.Errorf("Expected 3 genres, got %v", len(answer.Genre))
	}

	movieGenreTruth := []string{"Horror", "Mystery", "Thriller"}

	movieGenreAnswer, err := queries.GetGenreNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)

	if err != nil {
		t.Errorf("Encountered error querying for movie genre: %v", err)
	}

	if !cmp.Equal(movieGenreTruth, movieGenreAnswer) {
		t.Errorf("Expected %v, got %v", movieGenreTruth, movieGenreAnswer)
	}

	// Movie actor rows.
	if len(answer.Actor) != 3 {
		t.Errorf("Expected 3 actors, got %v", len(answer.Actor))
	}
	movieActorTruth := []string{"Anthony Franciosa", "Giuliano Gemma", "John Saxon"}
	movieActorAnswer, err := queries.GetActorNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying movie_actor: %v", err)
	}
	if !cmp.Equal(movieActorTruth, movieActorAnswer) {
		t.Errorf("Expected %v, got %v", movieActorTruth, movieActorAnswer)
	}

	// Movie director rows.
	if len(answer.Director) != 1 {
		t.Errorf("Expected 1 director uuid, got %v", len(answer.Director))
	}
	movieDirectorTruth := []string{"Dario Argento"}
	movieDirectorAnswer, err := queries.GetDirectorNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie director: %v", err)
	}
	if !cmp.Equal(movieDirectorTruth, movieDirectorAnswer) {
		t.Errorf("Expected %v, got %v", movieDirectorTruth, movieDirectorAnswer)
	}

	// Movie writer rows.
	if len(answer.Writer) != 1 {
		t.Errorf("Expected 1 writer uuid, got %v", len(answer.Writer))
	}
	movieWriterTruth := []string{"Dario Argento"}
	movieWriterAnswer, err := queries.GetWriterNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie writer: %v", err)
	}
	if !cmp.Equal(movieWriterTruth, movieWriterAnswer) {
		t.Errorf("Expected %v, got %v", movieWriterTruth, movieWriterAnswer)
	}
	// Movie rating rows.
	if len(answer.Rating) != 3 {
		t.Errorf("Expected 3 rating uuids, got %v", len(answer.Rating))
	}
	movieRatingAnswer, err := queries.GetRatingsForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie ratings: %v", err)
	}
	if len(movieRatingAnswer) != 3 {
		t.Errorf("Expected 3 rating answers, got %v", len(movieRatingAnswer))
	}
	movieRatingTruth := []database.MovieRating{
		{
			Uuid:            answer.Rating[0],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Internet Movie Database",
			Value:           "7.0/10.0",
			CreatedDatetime: movieRatingAnswer[0].CreatedDatetime,
		}, {
			Uuid:            answer.Rating[1],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Rotten Tomatoes",
			Value:           "77%",
			CreatedDatetime: movieRatingAnswer[1].CreatedDatetime,
		}, {
			Uuid:            answer.Rating[2],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Metacritic",
			Value:           "83/100",
			CreatedDatetime: movieRatingAnswer[2].CreatedDatetime,
		},
	}

	if !cmp.Equal(movieRatingTruth, movieRatingAnswer) {
		t.Errorf("Expected \n%v, got \n%v", movieRatingTruth, movieRatingAnswer)
	}
}
