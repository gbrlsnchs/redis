package redis

type Row struct {
	err  error
	rows *Rows
}

func (r *Row) Scan(value interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.rows.Next() {
		if !r.rows.multi {
			defer r.rows.Close()
		}

		return r.rows.Scan(value)
	}

	return r.rows.Close()
}
