// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testMovieGenres(t *testing.T) {
	t.Parallel()

	query := MovieGenres()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testMovieGenresDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMovieGenresQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := MovieGenres().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMovieGenresSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MovieGenreSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMovieGenresExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := MovieGenreExists(ctx, tx, o.UUID)
	if err != nil {
		t.Errorf("Unable to check if MovieGenre exists: %s", err)
	}
	if !e {
		t.Errorf("Expected MovieGenreExists to return true, but got false.")
	}
}

func testMovieGenresFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	movieGenreFound, err := FindMovieGenre(ctx, tx, o.UUID)
	if err != nil {
		t.Error(err)
	}

	if movieGenreFound == nil {
		t.Error("want a record, got nil")
	}
}

func testMovieGenresBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = MovieGenres().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testMovieGenresOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := MovieGenres().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testMovieGenresAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	movieGenreOne := &MovieGenre{}
	movieGenreTwo := &MovieGenre{}
	if err = randomize.Struct(seed, movieGenreOne, movieGenreDBTypes, false, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}
	if err = randomize.Struct(seed, movieGenreTwo, movieGenreDBTypes, false, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = movieGenreOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = movieGenreTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MovieGenres().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testMovieGenresCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	movieGenreOne := &MovieGenre{}
	movieGenreTwo := &MovieGenre{}
	if err = randomize.Struct(seed, movieGenreOne, movieGenreDBTypes, false, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}
	if err = randomize.Struct(seed, movieGenreTwo, movieGenreDBTypes, false, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = movieGenreOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = movieGenreTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testMovieGenresInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMovieGenresInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(movieGenreColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMovieGenreToOneMovieUsingMovie(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local MovieGenre
	var foreign Movie

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, movieDBTypes, true, movieColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Movie struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.MovieUUID, foreign.UUID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Movie().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.UUID, foreign.UUID) {
		t.Errorf("want: %v, got %v", foreign.UUID, check.UUID)
	}

	slice := MovieGenreSlice{&local}
	if err = local.L.LoadMovie(ctx, tx, false, (*[]*MovieGenre)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Movie == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Movie = nil
	if err = local.L.LoadMovie(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Movie == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testMovieGenreToOneSetOpMovieUsingMovie(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a MovieGenre
	var b, c Movie

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, movieGenreDBTypes, false, strmangle.SetComplement(movieGenrePrimaryKeyColumns, movieGenreColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, movieDBTypes, false, strmangle.SetComplement(moviePrimaryKeyColumns, movieColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, movieDBTypes, false, strmangle.SetComplement(moviePrimaryKeyColumns, movieColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Movie{&b, &c} {
		err = a.SetMovie(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Movie != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.MovieGenres[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.MovieUUID, x.UUID) {
			t.Error("foreign key was wrong value", a.MovieUUID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.MovieUUID))
		reflect.Indirect(reflect.ValueOf(&a.MovieUUID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.MovieUUID, x.UUID) {
			t.Error("foreign key was wrong value", a.MovieUUID, x.UUID)
		}
	}
}

func testMovieGenreToOneRemoveOpMovieUsingMovie(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a MovieGenre
	var b Movie

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, movieGenreDBTypes, false, strmangle.SetComplement(movieGenrePrimaryKeyColumns, movieGenreColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, movieDBTypes, false, strmangle.SetComplement(moviePrimaryKeyColumns, movieColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetMovie(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveMovie(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Movie().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Movie != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.MovieUUID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.MovieGenres) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testMovieGenresReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMovieGenresReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MovieGenreSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMovieGenresSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MovieGenres().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	movieGenreDBTypes = map[string]string{`UUID`: `TEXT`, `MovieUUID`: `TEXT`, `Name`: `TEXT`, `CreatedDatetime`: `INTEGER`}
	_                 = bytes.MinRead
)

func testMovieGenresUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(movieGenrePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(movieGenreAllColumns) == len(movieGenrePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenrePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testMovieGenresSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(movieGenreAllColumns) == len(movieGenrePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MovieGenre{}
	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenreColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, movieGenreDBTypes, true, movieGenrePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(movieGenreAllColumns, movieGenrePrimaryKeyColumns) {
		fields = movieGenreAllColumns
	} else {
		fields = strmangle.SetComplement(
			movieGenreAllColumns,
			movieGenrePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := MovieGenreSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testMovieGenresUpsert(t *testing.T) {
	t.Parallel()
	if len(movieGenreAllColumns) == len(movieGenrePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := MovieGenre{}
	if err = randomize.Struct(seed, &o, movieGenreDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MovieGenre: %s", err)
	}

	count, err := MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, movieGenreDBTypes, false, movieGenrePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MovieGenre struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MovieGenre: %s", err)
	}

	count, err = MovieGenres().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
