package sos_post

type Condition struct {
	ID        int    `filed:"id"`
	Name      string `filed:"name"`
	CreatedAt string `filed:"created_at"`
	UpdatedAt string `filed:"update_at"`
	DeletedAt string `filed:"deleted_at"`
}

type ConditionView struct {
	ID   int    `filed:"id"`
	Name string `filed:"name"`
}
