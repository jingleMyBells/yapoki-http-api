package internal

import (
	"log"
	"net/http"
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("lol")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user_id, err := GetUserIdByCookie(authCookie.Value)
		if err != nil {
			log.Printf("Ошибка поиска пользователя в базе по куки: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user_id == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
        next.ServeHTTP(w, r)
    })
}