package database

import (
	"github.com/pet-sitter/pets-next-door-api/internal/models"
	"github.com/pet-sitter/pets-next-door-api/internal/views"
)

func (tx *Tx) CreateUser(request *views.RegisterUserRequest) (*models.User, error) {
	user := &models.User{}

	err := tx.QueryRow(`
	INSERT INTO
		users
		(
			email,
			nickname,
			fullname,
			password,
			profile_image_id,
			fb_provider_type,
			fb_uid,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	RETURNING id, email, nickname, fullname, profile_image_id, fb_provider_type, fb_uid, created_at, updated_at
	`,
		request.Email,
		request.Nickname,
		request.Fullname,
		"",
		request.ProfileImageID,
		request.FirebaseProviderType,
		request.FirebaseUID,
	).Scan(&user.ID, &user.Email, &user.Nickname, &user.Fullname, &user.ProfileImageID, &user.FirebaseProviderType, &user.FirebaseUID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (tx *Tx) FindUserByEmail(email string) (*models.UserWithProfileImage, error) {
	user := &models.UserWithProfileImage{}

	err := tx.QueryRow(`
	SELECT
		users.id,
		users.email,
		users.nickname,
		users.fullname,
		media.url AS profile_image_url,
		users.fb_provider_type,
		users.fb_uid,
		users.created_at,
		users.updated_at
	FROM
		users
	LEFT OUTER JOIN
		media
	ON
		users.profile_image_id = media.id
	WHERE
		users.email = $1
	`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageURL,
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

func (tx *Tx) FindUserByUID(uid string) (*models.UserWithProfileImage, error) {
	user := &models.UserWithProfileImage{}

	err := tx.QueryRow(`
	SELECT
		users.id,
		users.email,
		users.nickname,
		users.fullname,
		media.url AS profile_image_url,
		users.fb_provider_type,
		users.fb_uid,
		users.created_at,
		users.updated_at
	FROM
		users
	LEFT JOIN
		media
	ON
		users.profile_image_id = media.id
	WHERE
		users.fb_uid = $1
	`,
		uid,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageURL,
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

func (tx *Tx) FindUserStatusByEmail(email string) (*models.UserStatus, error) {
	var userStatus models.UserStatus

	err := tx.QueryRow(`
	SELECT
		fb_provider_type
	FROM
		users
	WHERE
		email = $1
	`,
		email,
	).Scan(
		&userStatus.FirebaseProviderType,
	)

	if err != nil {
		return nil, err
	}

	return &userStatus, nil
}

func (tx *Tx) UpdateUserByUID(uid string, nickname string, profileImageID int) (*models.User, error) {
	user := &models.User{}

	err := tx.QueryRow(`
	UPDATE
		users
	SET
		nickname = $1,
		profile_image_id = $2,
		updated_at = NOW()
	WHERE
		fb_uid = $3
	RETURNING id, email, nickname, fullname, profile_image_id, fb_provider_type, fb_uid, created_at, updated_at
	`,
		nickname,
		profileImageID,
		uid,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageID,
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
