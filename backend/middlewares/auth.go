package middlewares

import (
	"net/http"

	"umigame-api/myerrors"
	"umigame-api/services"
	"umigame-api/utils"

	"github.com/jmoiron/sqlx"
)

type Auth struct {
	db      *sqlx.DB
	service *services.Service
}

func NewAuthMiddlewarer(db *sqlx.DB, service *services.Service) *Auth {
	return &Auth{
		db:      db,
		service: service,
	}
}

func (a *Auth) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		uuid, err := req.Cookie("uuid")
		if err != nil {
			myerrors.BadCookie.Wrap(err, "please set uuid in cookie")
			myerrors.ErrorHandler(w, req, err)
			return
		}

		userID, err := authorize(a.db, uuid.Value)
		if err != nil {
			myerrors.InvalidUUID.Wrap(ErrNoData, "uuid is not correct")
			myerrors.ErrorHandler(w, req, err)
			return
		}

		a.service.SetUserID(userID)

		cookie, err := utils.GetCookie(uuid.Value)
		if err != nil {
			myerrors.GenCookieFailed.Wrap(ErrNoData, "internal server error")
			myerrors.ErrorHandler(w, req, err)
		}

		http.SetCookie(w, cookie)

		next.ServeHTTP(w, req)
	})
}

func authorize(db *sqlx.DB, uuid string) (int, error) {
	sqlStr := `
	SELECT id
	FROM users
	WHERE uuid = ?;
	`

	var userID int
	if err := db.QueryRowx(sqlStr, uuid).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}
