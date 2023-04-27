package middlewares

//func AuthMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		var token string
//		cookie, err := r.Cookie("token")
//
//		authorizationHeader := r.Header.Get("Authorization")
//		fields := strings.Fields(authorizationHeader)
//
//		if len(fields) != 0 && fields[0] == "Bearer" {
//			token = fields[1]
//		} else if err == nil {
//			token = cookie.Value
//		}
//
//		if token == "" {
//			utils.JSONError(w, http.StatusUnauthorized, errors.New("you are not logged in"))
//			return
//		}
//
//		sub, err := utils.ValidateToken(token, config.Auth.JwtSecret)
//		if err != nil {
//			utils.JSONError(w, http.StatusUnauthorized, err)
//			return
//		}
//
//		user, err := storage.GetUserByID(context.Background(), int(sub.(float64)))
//		if err != nil {
//			utils.JSONError(w, http.StatusBadRequest, err)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "currentUser", user)
//
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
