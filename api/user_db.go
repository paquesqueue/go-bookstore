package api

func (db Query) InsertUser(req RequestUser) (ResponseUser, error) {
	const query = `INSERT INTO users 
	(username, email, fullname, hashed_password) 
	VALUES ($1, $2, $3, $4) 
	RETURNING username, email, fullname, hashed_password, created_at;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseUser{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Username, req.Email, req.Fullname, req.Password)

	resp := ResponseUser{}
	err = row.Scan(&resp.Username, &resp.Email, &resp.Fullname, &resp.HashedPassword, &resp.CreatedAt)
	if err != nil {
		return ResponseUser{}, err
	}
	return resp, nil
}

func (db Query) SelectUser(username string) (ResponseUser, error) {
	const query = `SELECT * FROM users WHERE username = $1;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseUser{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	resp := ResponseUser{}
	err = row.Scan(&resp.Username, &resp.Email, &resp.Fullname, &resp.HashedPassword, &resp.CreatedAt)
	if err != nil {
		return ResponseUser{}, err
	}
	return resp, nil
}

func (db Query) UpdateUser(username string, req RequestUser) (ResponseUser, error) {
	const query = `UPDATE users 
	SET username = $1, email = $2, fullname = $3, hashed_password = $4
	WHERE username = $5
	RETURNING username, email, fullname, hashed_password, created_at;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return ResponseUser{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Username, req.Email, req.Fullname, req.Password, username)
	resp := ResponseUser{}
	err = row.Scan(&resp.Username, &resp.Email, &resp.Fullname, &resp.HashedPassword, &resp.CreatedAt)
	if err != nil {
		return ResponseUser{}, err
	}
	return resp, nil
}

func (db Query) DeleteUser(username string) error {
	const query = `DELETE FROM users WHERE username = $1`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}
	return nil
}
