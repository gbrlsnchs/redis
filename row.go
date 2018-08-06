package redis

type Row struct {
	err  error
	rows *Rows
}

func (r *Row) Scan(value interface{}) error {
	if r.err != nil {
		defer r.rows.Close()
		return r.err
	}
	if r.rows.Next() {
		defer r.rows.Close()
		if err := r.rows.Err(); err != nil {
			return err
		}
		return r.rows.Scan(value)
	}
	return r.rows.Close()
}
