package main

import (
	"time"
	"strconv"
)

func XTIME_Get(T string) (*time.Time, error) {
	/*
	 * a little generic time generator:
	 *
	 * "" | nil = time.Now()
	 * 0 = beginning of time
	 * -offset = time.Now() minus offset
	 * YYYY-MM-DDTHH:MM:SS = exact time
	 */

	if T == "" || T == "nil" || T == "now" {
		t := time.Now().UTC()
		return &t, nil
	}

	if T == "0" {
		t, err := time.Parse("2006-01-02T15:04:05", "0000-01-01T00:00:00")
		if err != nil {
			return nil, err
		}
		t2 := t.UTC()
		return &t2, nil
	}

	if T[0] == '-' {
		/* Relative offset */
		t := time.Now().UTC()
		i, err := strconv.Atoi(T)
		if err != nil {
			return nil, err
		}
		t2 := t.Add(time.Duration(i)*time.Second)
		return &t2, nil
	}

	ts, err := time.Parse("2006-01-02T15:04:05", T)
	if err != nil {
		return nil, err
	}

	ts = ts.UTC()

	return &ts, nil
}


func LIMIT_Get(L string) (int, error) {
	limit := 50

	val, err := strconv.Atoi(L)
	if err != nil {
		return 0, err
	}

	if val <= 0 {
		limit = 100000000
	} else {
		limit = val
	}

	return limit, nil
}
