package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func CreateUser(ctx context.Context, tx database.Tx, request *user.RegisterUserRequest) (*user.User, *pnd.AppError) {
	const sql = `
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
	`

	user := &user.User{}
	if err := tx.QueryRowContext(ctx, sql,
		request.Email,
		request.Nickname,
		request.Fullname,
		"",
		request.ProfileImageID,
		request.FirebaseProviderType,
		request.FirebaseUID,
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
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return user, nil
}

func FindUsers(ctx context.Context, tx database.Tx, page int, size int, nickname *string) (*user.UserWithoutPrivateInfoList, *pnd.AppError) {
	const sql = `
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

	userList := user.NewUserWithoutPrivateInfoList(page, size)
	rows, err := tx.QueryContext(ctx, sql, nickname, size+1, (page-1)*size)
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

	userList.CalcLastPage()
	return userList, nil
}

func FindUserByID(ctx context.Context, tx database.Tx, id int, includeDeleted bool) (*user.UserWithProfileImage, *pnd.AppError) {
	const sql = `
	SELECT
		users.id,
		users.email,
		users.nickname,
		users.fullname,
		media.url AS profile_image_url,
		users.fb_provider_type,
		users.fb_uid,
		users.created_at,
		users.updated_at,
		users.deleted_at
	FROM
		users
	LEFT OUTER JOIN
		media
	ON
		users.profile_image_id = media.id
	WHERE
		users.id = $1 AND
		(users.deleted_at IS NULL OR $2)
	`

	var user user.UserWithProfileImage
	if err := tx.QueryRowContext(ctx, sql, id, includeDeleted).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageURL,
		&user.FirebaseProviderType,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &user, nil
}

func FindUserByEmail(ctx context.Context, tx database.Tx, email string) (*user.UserWithProfileImage, *pnd.AppError) {
	const sql = `
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
	`

	var user user.UserWithProfileImage
	if err := tx.QueryRowContext(ctx, sql, email).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageURL,
		&user.FirebaseProviderType,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &user, nil
}

func FindUserByUID(ctx context.Context, tx database.Tx, uid string) (*user.UserWithProfileImage, *pnd.AppError) {
	const sql = `
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
	`

	var user user.UserWithProfileImage
	if err := tx.QueryRowContext(ctx, sql, uid).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Fullname,
		&user.ProfileImageURL,
		&user.FirebaseProviderType,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &user, nil
}

func FindUserIDByFbUID(ctx context.Context, tx database.Tx, fbUid string) (int, *pnd.AppError) {
	const sql = `
	SELECT
		id
	FROM
		users
	WHERE
		fb_uid = $1 AND
		deleted_at IS NULL
	`

	var userID int
	if err := tx.QueryRowContext(ctx, sql, fbUid).Scan(&userID); err != nil {
		return 0, pnd.FromPostgresError(err)
	}

	return userID, nil
}

func ExistsUserByNickname(ctx context.Context, tx database.Tx, nickname string) (bool, *pnd.AppError) {
	const sql = `
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
	`

	var exists bool
	if err := tx.QueryRowContext(ctx, sql, nickname).Scan(&exists); err != nil {
		return false, pnd.FromPostgresError(err)
	}

	return exists, nil
}

func FindUserStatusByEmail(ctx context.Context, tx database.Tx, email string) (*user.UserStatus, *pnd.AppError) {
	const sql = `
	SELECT
		fb_provider_type
	FROM
		users
	WHERE
		email = $1 AND
		deleted_at IS NULL
	`

	var userStatus user.UserStatus
	if err := tx.QueryRowContext(ctx, sql, email).Scan(&userStatus.FirebaseProviderType); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &userStatus, nil
}

func UpdateUserByUID(ctx context.Context, tx database.Tx, uid string, nickname string, profileImageID *int) (*user.User, *pnd.AppError) {
	const sql = `
	UPDATE
		users
	SET
		nickname = $1,
		profile_image_id = $2,
		updated_at = NOW()
	WHERE
		fb_uid = $3 AND
		deleted_at IS NULL
	RETURNING
		id,
		email,
		nickname,
		fullname,
		profile_image_id,
		fb_provider_type,
		fb_uid,
		created_at,
		updated_at
	`

	var user user.User
	err := tx.QueryRowContext(ctx, sql,
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

	return &user, nil
}

func DeleteUserByUID(ctx context.Context, tx database.Tx, uid string) *pnd.AppError {
	const sql = `
	UPDATE
		users
	SET
		deleted_at = NOW()
	WHERE
		fb_uid = $1
	`

	if _, err := tx.ExecContext(ctx, sql, uid); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}
