package internal

import (
	"time"
)

type Variant struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Problem struct {
	UserId int `json:"id"`
	VariantId int `json:"variant_id"`
	Question string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Answer1 string `json:"answer_1"`
	Answer2 string `json:"answer_2"`
	Answer3 string `json:"answer_3"`
}

type ProblemOutput struct {
	Question string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Answer1 string `json:"answer_1"`
	Answer2 string `json:"answer_2"`
	Answer3 string `json:"answer_3"`
}

type Test struct {
	UserId int `json:"id"`
	VariantId int `json:"variant_id"`
	StartTime time.Time `json:"start_time"`
}

type TestAnswer struct {
	TestId int `json:"test_id"`
	ProblemId int `json:"problem_id"`
	Answer string `json:"answer"`
}

type TestResult struct {
	TestId int `json:"test_id"`
	Percent int `json:"percent"`
}

type User struct {
	Login string `json:"login"`
	Password string `json:"password"`
	LoginTime time.Time `json:"login_time"`
	LogoutTime time.Time `json:"logout_time"`
	Cookie string `json:"cookie"`
}