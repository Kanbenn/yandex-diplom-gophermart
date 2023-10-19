package models

type UserInsert struct {
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
	UID         int
	WasConflict bool
	Err         error
}

// had to make the static-check happy >_<
type CtxKey string

const CtxKeyUser CtxKey = "uid"
