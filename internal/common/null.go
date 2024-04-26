package utils

import "database/sql"

func DerefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}

func IsNotNil[T any](val *T) bool {
	return val != nil
}

func NullStrToStrPtr(val sql.NullString) *string {
	if val.Valid {
		return &val.String
	}
	return nil
}

func StrPtrToNullStr(val *string) sql.NullString {
	return sql.NullString{
		String: DerefOrEmpty(val),
		Valid:  IsNotNil(val),
	}
}

func StrToNullStr(val string) sql.NullString {
	return sql.NullString{
		String: val,
		Valid:  val != "",
	}
}

func NullInt64ToInt64Ptr(val sql.NullInt64) *int64 {
	if val.Valid {
		return &val.Int64
	}
	return nil
}

func IntToNullInt64(val int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(val),
		Valid: val != 0,
	}
}

func Int64PtrToNullInt64(val *int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: DerefOrEmpty(val),
		Valid: IsNotNil(val),
	}
}

func IntPtrToNullInt64(val *int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(DerefOrEmpty(val)),
		Valid: IsNotNil(val),
	}
}

func IntToNullInt32(val int) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(val),
		Valid: val != 0,
	}
}

func Int64ToNullInt32(val int64) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(val),
		Valid: val != 0,
	}
}

func IntPtrToNullInt32(val *int) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(DerefOrEmpty(val)),
		Valid: IsNotNil(val),
	}
}

func Int64PtrToNullInt32(val *int64) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(DerefOrEmpty(val)),
		Valid: IsNotNil(val),
	}
}
