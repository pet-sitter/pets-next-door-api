package postgres

import (
	"github.com/pet-sitter/pets-next-door-api/internal/database"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type UserPostgresStore struct {
	db *database.DB
}

func NewUserPostgresStore(db *database.DB) *UserPostgresStore {
	return &UserPostgresStore{
		db: db,
	}
}

func (s *UserPostgresStore) CreateUser(request *user.RegisterUserRequest) (*user.User, error) {
	user := &user.User{}

	tx, _ := s.db.Begin()

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
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserByEmail(email string) (*user.UserWithProfileImage, error) {
	user := &user.UserWithProfileImage{}

	tx, _ := s.db.Begin()
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
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserByUID(uid string) (*user.UserWithProfileImage, error) {
	user := &user.UserWithProfileImage{}

	tx, _ := s.db.Begin()
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
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserStatusByEmail(email string) (*user.UserStatus, error) {
	var userStatus user.UserStatus

	tx, _ := s.db.Begin()
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
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return &userStatus, nil
}

func (s *UserPostgresStore) UpdateUserByUID(uid string, nickname string, profileImageID int) (*user.User, error) {
	user := &user.User{}

	tx, _ := s.db.Begin()
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
	tx.Commit()

	if err != nil {
		return nil, err
	}

	return user, nil
}
