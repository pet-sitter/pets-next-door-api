package postgres

import (
	"context"

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

func (s *UserPostgresStore) CreateUser(ctx context.Context, request *user.RegisterUserRequest) (*user.User, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	user := &user.User{}
	err = tx.QueryRowContext(ctx, `
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
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUsers(ctx context.Context, page int, size int, nickname *string) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	userList := user.NewUserWithoutPrivateInfoList(page, size)

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

	rows, err := tx.QueryContext(ctx, query, nickname, size+1, (page-1)*size)
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
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	userList.CalcLastPage()
	return userList, nil
}

func (s *UserPostgresStore) FindUserByEmail(ctx context.Context, email string) (*user.UserWithProfileImage, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	user := &user.UserWithProfileImage{}
	err = tx.QueryRowContext(ctx, `
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
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserByUID(ctx context.Context, uid string) (*user.UserWithProfileImage, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	user := &user.UserWithProfileImage{}
	err = tx.QueryRowContext(ctx, `
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
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func (s *UserPostgresStore) FindUserIDByFbUID(ctx context.Context, fbUid string) (int, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return 0, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	var UserID int
	err = tx.QueryRowContext(ctx, `
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
	if err != nil {
		return 0, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return 0, pnd.FromPostgresError(err)
	}

	return UserID, nil
}

func (s *UserPostgresStore) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `
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

	if err := tx.Commit(); err != nil {
		return false, pnd.FromPostgresError(err)
	}

	return exists, nil
}

func (s *UserPostgresStore) FindUserStatusByEmail(ctx context.Context, email string) (*user.UserStatus, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	var userStatus user.UserStatus
	err = tx.QueryRowContext(ctx, `
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
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &userStatus, nil
}

func (s *UserPostgresStore) UpdateUserByUID(ctx context.Context, uid string, nickname string, profileImageID *int) (*user.User, *pnd.AppError) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer tx.Rollback()

	user := &user.User{}
	err = tx.QueryRowContext(ctx, `
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
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}
