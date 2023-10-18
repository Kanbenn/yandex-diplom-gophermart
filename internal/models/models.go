package models

type User struct {
	ID       int    `json:"-"` // hide from json
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderInsert struct {
	Number     string
	UID        int
	UploadedAt string
}
type OrderInsertResult struct {
	UID        int
	IsConflict bool
	Err        error
}

// had to make the static-check happy >_<
type CtxKey string

const CtxKeyUser CtxKey = "uid"
