package internal

type Result struct {
	last int
}

func (r *Result) LastInsertId() (int64, error) {
	return int64(r.last), nil
}

func (r *Result) RowsAffected() (int64, error) {
	return 0, nil
}
