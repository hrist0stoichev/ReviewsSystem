package middlewares

import (
	"net/http"
	"strings"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
)

const (
	UserIDHeader                   = "ReviewsSystem-UserID"
	UserRoleHeader                 = "ReviewsSystem-UserRole"
	authorizationHeader            = "Authorization"
	bearerTokenPrefix              = "Bearer"
	forbiddenErrorMessage          = "Forbidden"
	unauthorizedErrorMessage       = "Unauthorized"
	unsupportedAuthorizationMethod = "Unsupported Authorization Method"
)

type authMiddleware struct {
	tokensService services.TokensService
	logger        log.Logger
}

func NewAuth(tokensService services.TokensService) *authMiddleware {
	return &authMiddleware{
		tokensService: tokensService,
	}
}

func (ah *authMiddleware) AuthorizeForRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqToken := r.Header.Get(authorizationHeader)
			if reqToken == "" {
				http.Error(w, unauthorizedErrorMessage, http.StatusUnauthorized)
				return
			}

			splitToken := strings.Split(reqToken, " ")
			if len(splitToken) != 2 {
				http.Error(w, unauthorizedErrorMessage, http.StatusUnauthorized)
				return
			}

			switch splitToken[0] {
			case bearerTokenPrefix:
				userClaims, err := ah.tokensService.ParseSignedToken(splitToken[1])
				if err != nil {
					ah.logger.WithError(err).Warnln("could not parse signed token")
					http.Error(w, unauthorizedErrorMessage, http.StatusUnauthorized)
					return
				}

				if !contains(userClaims.Role, roles...) {
					http.Error(w, forbiddenErrorMessage, http.StatusForbidden)
					return
				}

				r.Header.Add(UserIDHeader, userClaims.Id)
				r.Header.Add(UserRoleHeader, userClaims.Role)
			default:
				http.Error(w, unsupportedAuthorizationMethod, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserUIDFromRequest returns the user ID from an incoming request (empty string if not authenticated)
func UserIDFromRequest(r *http.Request) string {
	return r.Header.Get(UserIDHeader)
}

func contains(str string, strs ...string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
