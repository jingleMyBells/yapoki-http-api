package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)


func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Root handled sucessfully")
}

func PathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "/lalala handled sucessfully")
}


func VariantsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			variants := GetAllVariants()
			json.NewEncoder(w).Encode(variants)
			return
		case http.MethodPost:
			pathParts := strings.Split((r.URL.Path), "/")

			if len(pathParts) < 3 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный URL"))
				return
			}

			variantNumber, err := strconv.Atoi(pathParts[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
			}

			if variantNumber == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Невалидный вариант"))
			}

			// тут ндо пофиксить когда появится авторизация
			user_id := 1

			test, err := CreateTest(variantNumber, user_id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			fmt.Println(test)

			pathString := "/testing/" + strconv.Itoa(test.VariantId)

			http.Redirect(w, r, pathString, http.StatusBadRequest)

	}
}

func TestingHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split((r.URL.Path), "/")

	// надо пофиксить тут когда появится авторизация
	user_id := 1

	if len(pathParts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный URL"))
		return
	}

	if len(pathParts) == 3 {
		// ветвление про то, что мы еще не знаем какие задачи решаем, только вариант
		// надо провалидировать что это инт
		variantId, err := strconv.Atoi(pathParts[2])
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

		lastUserTest, err := GetTLastTestIdByVariantAndUser(variantId, user_id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Тестирование не найдено"))
			return
		}
		
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
				pathString := "/testing/result/" + strconv.Itoa(lastUserTest)
				http.Redirect(w, r, pathString, http.StatusBadRequest)
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(availableProblems)
		return

	} else if len(pathParts) > 3 {
		// ветвление про то, когда мы знаем И номер варианта, И номер задачи
		// надо провалидировать что это инт
		variantId, err := strconv.Atoi(pathParts[2])
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
		
		problemId, err := strconv.Atoi(pathParts[3])
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
			problem ProblemOutput
			availableProblems []int
		}

		lastUserTest, err := GetTLastTestIdByVariantAndUser(variantId, user_id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Тестирование не найдено"))
			return
		}
		
		availableProblems, err := GetTestUnsolvedProblems(lastUserTest)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Задачи не найдены"))
			return
		}

		
		response := responseStruct {
			problem: *problem,
			availableProblems: availableProblems,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split((r.URL.Path), "/")

	if len(pathParts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Невалидный URL"))
		return
	}

	testId, err := strconv.Atoi(pathParts[2])
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResult)
	return
}


