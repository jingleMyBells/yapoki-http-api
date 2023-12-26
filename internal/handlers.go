package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "../../website/variants.html")
}

func TestHTMLHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../website/testing.html")
}

func TestResultHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../website/results.html")
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../website/login.html")
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		type reqBody struct {
			Login string  `json:"login"`
			Password string `json:"password"`
		}
		var body reqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("Ошибка чтения тела запроса: %v", err)
		}

		user_id, cookie, err := GetUserByLoginPassword(body.Login, body.Password)
		if err != nil {
			log.Printf("Не удалось получить последний тест юзера: %v", err)
		}
		if user_id == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{'Хендлер: Юзер не найден"))
			return
		}
		if cookie == "" {
			cookie, err = CreateUserCookie(user_id)
			if err != nil {
				log.Printf("Ошибка записи куки в базу: %v", err)
			}
			cookieToInstall := http.Cookie{
				Name:   GetAppConfig().CookieName,
				Value:  cookie,
				MaxAge: 86400,
				SameSite: http.SameSiteNoneMode,
				Secure: true,
				Expires: time.Now().Add(365 * 24 * time.Hour),
				HttpOnly: true,
				Path: "/",
			}
			http.SetCookie(w, &cookieToInstall)
		} else {
			cookieToInstall := http.Cookie{
				Name:   GetAppConfig().CookieName,
				Value:  cookie,
				MaxAge: 86400,
				SameSite: http.SameSiteNoneMode,
				Secure: true,
				Expires: time.Now().Add(365 * 24 * time.Hour),
				HttpOnly: true,
				Path: "/",
			}
			http.SetCookie(w, &cookieToInstall)
		}
		
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не разрешен"))
		return
	}

}

func VariantsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		variants := GetAllVariants()
		w.Header().Set("Content-Type", "application/json; charset: utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(variants)
		return
	case http.MethodPost:
		pathParts := strings.Split((r.URL.Path), "/")

		if len(pathParts) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Невалидный URL"))
			return
		}

		variantNumber, err := strconv.Atoi(pathParts[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Невалидный вариант"))
			return
		}

		if variantNumber == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Невалидный вариант"))
			return
		}

		authCookie, err := r.Cookie(GetAppConfig().CookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user_id, err := GetUserIdByCookie(authCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		test, err := CreateTest(variantNumber, user_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("Начато прохождение тестирования: %v", test)

		w.WriteHeader(http.StatusCreated)
		return
	}
}

func TestingHandler(w http.ResponseWriter, r *http.Request) {

	pathParts := strings.Split((r.URL.Path), "/")

	if len(pathParts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный URL"))
		return
	}

	variantId, err := strconv.Atoi(pathParts[3])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
				return
			}

			if variantId == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
				return
			}

	authCookie, err := r.Cookie(GetAppConfig().CookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user_id, err := GetUserIdByCookie(authCookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	lastUserTest, err := GetTLastTestIdByVariantAndUser(variantId, user_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Тестирование не найдено"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		
		if len(pathParts) == 4 {

			availableProblems, err := GetTestUnsolvedProblems(lastUserTest)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Задачи не найдены"))
				return
			}

			if len(availableProblems) == 0 {
				testFinished, err := CheckTestIsFinished(lastUserTest)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Не удалось проверить статус теста"))
					return
				}
				if testFinished {
					err = CreateTestResults(lastUserTest)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Не удалось проверить статус теста"))
						return
					}
					type responseStruct struct {
						IsFinished bool `json:"is_finished"`
						TestId int `json:"test_id"`
					}
					response := responseStruct{
						IsFinished: testFinished,
						TestId: lastUserTest,
					}
					w.Header().Set("Content-Type", "application/json; charset: utf-8")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json; charset: utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(availableProblems)
			return

		} else if len(pathParts) > 4 {
			variantId, err := strconv.Atoi(pathParts[3])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
				return
			}

			if variantId == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
				return
			}

			problemId, err := strconv.Atoi(pathParts[4])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный номер задачи"))
				return
			}

			if problemId == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный номер задачи"))
				return
			}

			problem, err := GetProblemById(problemId)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Задача не найдена"))
				return
			}

			type responseStruct struct {
				Problem           ProblemOutput `json:"problem"`
				AvailableProblems []int `json:"available_problems"`
			}

			availableProblems, err := GetTestUnsolvedProblems(lastUserTest)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Задачи не найдены"))
				return
			}

			response := responseStruct{
				Problem:           *problem,
				AvailableProblems: availableProblems,
			}

			w.Header().Set("Content-Type", "application/json; charset: utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	case http.MethodPost:
		type reqBody struct {
			Answer string  `json:"answer"`
			Problem_id string `json:"problem_id"`
		}
		var body reqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("Ошибка чтения тела запроса: %v", err)
		}
		problem_id, err := strconv.Atoi(body.Problem_id)
		if err != nil {
			log.Printf("Ошибка конвертации входящих данных: %v", err)
		}
		err = CreateTestAnswer(lastUserTest, problem_id, body.Answer)
		if err != nil {
			log.Printf("Ошибка создания тестироваания: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split((r.URL.Path), "/")

	if len(pathParts) < 5 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный URL"))
		return
	}

	testId, err := strconv.Atoi(pathParts[4])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный идентификатор теста"))
		return
	}

	if testId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный идентификатор теста"))
		return
	}

	testResult, err := GetTestResults(testId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Не удалось получить результаты теста"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset: utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResult)
	return
}
