package internal

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)


func GetUser(user_id int) error {
	db := GetDB()

	row := db.Sql.QueryRow("SELECT id, login FROM users WHERE id = $1", user_id) 
	if err := row.Scan(); err == sql.ErrNoRows {
		return err
	}

	return nil
}


func GetUserByLoginPassword(login string, pass string) (int, string, error) {

	inputPassword := []byte(pass)
	inputSha1Hash := fmt.Sprintf("%x", sha1.Sum(inputPassword))

	db := GetDB()
	row := db.Sql.QueryRow(`SELECT id, password, COALESCE(cookie, '') FROM users WHERE login = $1`, login) 
	var user_id int
	var cookie string
	var dbPassword string
	if err := row.Scan(&user_id, &dbPassword, &cookie); err == sql.ErrNoRows {
		log.Printf("Юзер не найден: %v", err)
		return 0, "", err
	} else if err != nil {
		log.Printf("Ошибка при попытке найти пользователя: %v", err)
		return 0, "", err
	}

	if inputSha1Hash == dbPassword {
		return user_id, cookie, nil
	}

	return 0, "", fmt.Errorf("Неверный логин или пароль")
}


func CreateUserCookie(user_id int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, GetAppConfig().CookieLength)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

	cookie := string(b)

	db := GetDB()
	_, err := db.Sql.Exec(`UPDATE users SET cookie = $1 WHERE id = $2`, cookie, user_id) 
	if err != nil {
		log.Printf("Ошибка записи в базу данных: %v", err)
		return "", err
	}

	return cookie, nil
}


func GetUserIdByCookie(cookie string) (int, error) {
	db := GetDB()

	row := db.Sql.QueryRow("SELECT id FROM users WHERE cookie = $1", cookie) 
	var user_id int
	if err := row.Scan(&user_id); err == sql.ErrNoRows {
		return 0, err
	} else if err != nil {
		return 0, err
	}

	return user_id, nil
}


func GetAllVariants() []Variant {
	db := GetDB()

	variants := make([]Variant, 0)

	rows, err := db.Sql.Query("SELECT id, name FROM variant")
	if err != nil {
		log.Printf("Ошибка чтения из базы: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		variant := Variant{
			Id: id,
			Name: name,
		}
		variants = append(variants, variant)
	}

	return variants
}


func CreateTest(variant_id int, user_id int) (*Test, error) {
	db := GetDB()

	test := Test{
		UserId: user_id,
		VariantId: variant_id,
		StartTime: time.Now(),
	}

	err := db.AddTest(test)
	if err != nil {
		log.Printf("Ошибка создания тестирования в базе: %v", err)
		return nil, err
	}

	return &test, nil
}


func GetTLastTestIdByVariantAndUser(variant_id int, user_id int) (int, error) {
	db := GetDB()

	test_row := db.Sql.QueryRow(`SELECT id 
	FROM test 
	WHERE variant_id = $1 AND user_id = $2 
	ORDER BY start_time DESC LIMIT 1`, variant_id, user_id)
	var test_id int
	if err := test_row.Scan(&test_id); err == sql.ErrNoRows {
		log.Printf("Не удалось получить последний тест юзера: %v", err)
		return 0, err
	}

	return test_id, nil
}


func GetTestUnsolvedProblems(test_id int) ([]int, error) {
	db := GetDB()

	rows, err := db.Sql.Query(`SELECT id
	FROM problem p 
	WHERE variant_id = (SELECT variant_id FROM test WHERE test.id = $1)
	EXCEPT
	SELECT problem_id
	FROM test_answer ta 
	WHERE ta.test_id = $1;`, test_id)
	if err != nil {
		log.Printf("Не удалось получить нерешенные задания: %v", err)
	}
	defer rows.Close()

	testUnresolvedProblems := make([]int, 0)

	for rows.Next() {
		var id int
		rows.Scan(&id)
		testUnresolvedProblems = append(testUnresolvedProblems, id)
	}

	return 	testUnresolvedProblems, nil
}


func GetProblemById(problem_id int) (*ProblemOutput, error) {
	db := GetDB()
	problem_row := db.Sql.QueryRow(`SELECT question, correct_answer, answer_1, answer_2, answer_3
	FROM problem WHERE id = $1`, problem_id)
	var question string
	var correctAnswer string
	var answer1 string
	var answer2 string
	var answer3 string
	if err := problem_row.Scan(&question, &correctAnswer, &answer1, &answer2, &answer3); err == sql.ErrNoRows {
		return nil, fmt.Errorf("Задание не найдено")
	} else if err != nil {
		return nil, err
	}

	problemOutput := ProblemOutput{
		Question: question,
		CorrectAnswer: correctAnswer,
		Answer1: answer1,
		Answer2: answer2,
		Answer3: answer3,
	}

	return &problemOutput, nil
}


func CreateTestAnswer(test_id int, problem_id int, answer string) error {
	db := GetDB()


	test_answer := TestAnswer{
		TestId: test_id,
		ProblemId: problem_id,
		Answer: answer,
	}

	err := db.AddTestAnswer(test_answer)
	if err != nil {
		return err
	}

	return nil
}


func CheckTestIsFinished(test_id int) (bool, error) {
	db := GetDB()

	variant_problems_rows, err := db.Sql.Query(`SELECT COUNT(*) 
	FROM problem 
	WHERE variant_id = (SELECT variant_id FROM test WHERE test.id = $1)`, test_id)
	defer variant_problems_rows.Close()
	var test_problems_count int
	if err != nil {
		return false, err
	}
	defer variant_problems_rows.Close()
	for variant_problems_rows.Next() {
		variant_problems_rows.Scan(&test_problems_count)
	}

	resolved_test_problems, err := db.Sql.Query(`SELECT COUNT(*)
	FROM test_answer WHERE test_id = $1`, test_id)
	var resolved_test_problems_count int
	if err != nil {
		return false, err
	}
	defer resolved_test_problems.Close()
	for resolved_test_problems.Next() {
		resolved_test_problems.Scan(&resolved_test_problems_count)
	}

	if test_problems_count <= resolved_test_problems_count {
		return true, nil
	}

	return false, nil
}


func CreateTestResults(test_id int) error {
	db := GetDB()

	correct_answers, err := db.Sql.Query(`SELECT COUNT(*)
	FROM test_answer ta 
	JOIN problem p ON p.id = ta.problem_id
	WHERE ta.test_id = $1 AND ta.answer = p.correct_answer 
	`, test_id)
	defer correct_answers.Close()
	var correct_answers_count int
	if err != nil {
		return err
	}
	for correct_answers.Next() {
		correct_answers.Scan(&correct_answers_count)
	}

	total_problems, err := db.Sql.Query(`SELECT COUNT(*)
	FROM problem p 
	WHERE variant_id = (SELECT variant_id FROM test WHERE test.id = $1)
	`, test_id)
	defer total_problems.Close()
	var total_problems_count int
	if err != nil {
		return err
	}
	for total_problems.Next() {
		total_problems.Scan(&total_problems_count)
	}

	correct_answers_count *= 100
	percent := correct_answers_count / total_problems_count

	testResult := TestResult{
		TestId: test_id,
		Percent: percent,
	}

	err = db.AddTestResult(testResult)
	if err != nil {
		return err
	}

	return nil
}


func GetTestResults(test_id int) (*TestResult, error) {
	db := GetDB()
	test_result_row := db.Sql.QueryRow(`SELECT percent
	FROM test_result WHERE test_id = $1`, test_id)
	var percent int
	if err := test_result_row.Scan(&percent); err == sql.ErrNoRows {
		return nil, fmt.Errorf("Вопрос не найден")
	} else if err != nil {
		return nil, err
	}
	testResult := TestResult{
		TestId: test_id,
		Percent: percent,
	}

	return &testResult, nil

}
