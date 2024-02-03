package postgres

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type UserPostgresStore struct {
	db *database.DB
}

func NewUserPostgresStore(db *database.DB) *UserPostgresStore {
	return &UserPostgresStore{
		db: db,
	}
}

func (s *UserPostgresStore) CreateUser(request *user.RegisterUserRequest) (*user.User, *pnd.AppError) {
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
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUsers(page int, size int, nickname *string) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	userList := user.NewUserWithoutPrivateInfoList(page, size)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	query := `
	SELECT
		users.id,
		users.nickname,
		media.url AS profile_image_url
	FROM
		users
	LEFT OUTER JOIN
		media
	ON
		users.profile_image_id = media.id
	WHERE
	    (users.nickname = $1 OR $1 IS NULL) AND
		users.deleted_at IS NULL
	ORDER BY
	    users.created_at DESC
	LIMIT $2
	OFFSET $3
	`

	rows, err := tx.Query(query, nickname, size+1, (page-1)*size)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &user.UserWithoutPrivateInfo{}

		err := rows.Scan(&user.ID, &user.Nickname, &user.ProfileImageURL)
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		userList.Items = append(userList.Items, *user)
	}
	userList.CalcLastPage()

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return userList, nil
}

func (s *UserPostgresStore) FindUserByEmail(email string) (*user.UserWithProfileImage, *pnd.AppError) {
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
		users.email = $1 AND
		users.deleted_at IS NULL
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
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserByUID(uid string) (*user.UserWithProfileImage, *pnd.AppError) {
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
		users.fb_uid = $1 AND
		users.deleted_at IS NULL
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
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserIDByFbUID(fbUid string) (int, *pnd.AppError) {
	var UserID int

	tx, _ := s.db.Begin()
	err := tx.QueryRow(`
	SELECT
		id
	FROM
		users
	WHERE
		fb_uid = $1 AND
		deleted_at IS NULL
	`,
		fbUid,
	).Scan(&UserID)
	tx.Commit()

	if err != nil {
		return 0, pnd.FromPostgresError(err)
	}

	return UserID, nil
}

func (s *UserPostgresStore) ExistsByNickname(nickname string) (bool, *pnd.AppError) {
	var exists bool

	tx, _ := s.db.Begin()
	err := tx.QueryRow(`
	SELECT
		CASE
		    WHEN EXISTS (
				SELECT
					1
				FROM
					users
				WHERE
					nickname = $1 AND
					deleted_at IS NULL
			) THEN TRUE
			ELSE FALSE
		END
	`,
		nickname,
	).Scan(&exists)
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}

	err = tx.Commit()
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}

	return exists, nil
}

func (s *UserPostgresStore) FindUserStatusByEmail(email string) (*user.UserStatus, *pnd.AppError) {
	var userStatus user.UserStatus

	tx, _ := s.db.Begin()
	err := tx.QueryRow(`
	SELECT
		fb_provider_type
	FROM
		users
	WHERE
		email = $1 AND
		deleted_at IS NULL
	`,
		email,
	).Scan(
		&userStatus.FirebaseProviderType,
	)
	tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &userStatus, nil
}

func (s *UserPostgresStore) UpdateUserByUID(uid string, nickname string, profileImageID *int) (*user.User, *pnd.AppError) {
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
		fb_uid = $3 AND
		deleted_at IS NULL
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
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}
