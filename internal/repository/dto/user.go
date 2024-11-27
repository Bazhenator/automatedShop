package dto

type User struct {
	ID       uint64 `db:"id"`
	Login    string `db:"login"`
	PassHash []byte `db:"pass_hash"`
}
