package database

import "github.com/pet-sitter/pets-next-door-api/internal/models"

func (tx *Tx) CreateUser(user *models.User) (*models.User, error) {
	err := tx.QueryRow(`
	INSERT INTO
		users
		(
			email,
			nickname,
			fullname,
			password,
			fb_provider_type,
			fb_uid,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`,
		user.Email,
		user.Nickname,
		user.Fullname,
		user.Password,
		user.FirebaseProviderType,
		user.FirebaseUID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (tx *Tx) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := tx.QueryRow(`
	SELECT
		id,
		email,
		nickname,
		fullname,
		fb_provider_type,
		fb_uid,
		created_at,
		updated_at
	FROM
		users
	WHERE
		email = $1
	`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.FirebaseProviderType,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (tx *Tx) FindUserByUID(uid string) (*models.User, error) {
	user := &models.User{}

	err := tx.QueryRow(`
	SELECT
		id,
		email,
		nickname,
		fullname,
		fb_provider_type,
		fb_uid,
		created_at,
		updated_at
	FROM
		users
	WHERE
		fb_uid = $1
	`,
		uid,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.FirebaseProviderType,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (tx *Tx) UpdateUserByUID(uid string, nickname string) (*models.User, error) {
	user := &models.User{}

	err := tx.QueryRow(`
	UPDATE
		users
	SET
		nickname = $1,
		updated_at = NOW()
	WHERE
		fb_uid = $2
	RETURNING id, created_at, updated_at
	`,
		nickname,
		uid,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
