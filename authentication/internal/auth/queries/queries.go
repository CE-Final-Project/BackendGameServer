package queries

type AuthQueries struct {
	Login       LoginAuthHandler
	VerifyToken VerifyAuthHandler
}
