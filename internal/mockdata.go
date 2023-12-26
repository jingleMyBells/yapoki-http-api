package internal

import (
	"time"
	"crypto/sha1"
	"fmt"
	"database/sql"
)




func AddMockData() error {
	db := GetDB()

	var userCount int
	row := db.Sql.QueryRow("SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err == sql.ErrNoRows {
		return err
	}

	fmt.Println(userCount)

	if userCount == 0 {
		inputPassword := []byte("12345")
		inputSha1Hash := fmt.Sprintf("%x", sha1.Sum(inputPassword))

		fmt.Println(inputSha1Hash)

		user := User{
			Login: "admin",
			Password: inputSha1Hash,
			LoginTime: time.Now(),
			LogoutTime: time.Now(),
		}

		variant1 := Variant{
			Name: "вариант1",
		}
		variant2 := Variant{
			Name: "вариант2",
		}
		variant3 := Variant{
			Name: "вариант2",
		}

		problem1 := Problem{
			VariantId: 1,
			Question: `Как назывался особый головной убор, 
			который носили фараоны в Древнем Египте?`,
			CorrectAnswer: "Немес",
			Answer1: "Картуз",
			Answer2: "Корона",
			Answer3: "Убрус",
		}
		problem2 := Problem{
			VariantId: 1,
			Question: "У какого животного самые большие глаза относительно тела?",
			CorrectAnswer: "У долгопята",
			Answer1: "У лемура",
			Answer2: "У летучей мыши",
			Answer3: "У тупайи",
		}
		problem3 := Problem{
			VariantId: 1,
			Question: "Какие огурцы сажал на брезентовом поле герой одноименной песни?",
			CorrectAnswer: "Алюминиевые",
			Answer1: "Оловянные",
			Answer2: "Медные",
			Answer3: "Железные",
		}

		problem4 := Problem{
			VariantId: 2,
			Question: "Какое из этих растений — плотоядное?",
			CorrectAnswer: "Росянка",
			Answer1: "Володушка",
			Answer2: "Мытник",
			Answer3: "Астрагал",
		}
		problem5 := Problem{
			VariantId: 2,
			Question: `Как называется логически верная ситуация, 
			которая не может существовать в реальности?`,
			CorrectAnswer: "Апория",
			Answer1: "Парадокс",
			Answer2: "Антиномия",
			Answer3: "Гипербола",
		}
		problem6 := Problem{
			VariantId: 2,
			Question: "С какой головой изображается индуистский бог Ганеша?",
			CorrectAnswer: "Слона",
			Answer1: "Обезьяны",
			Answer2: "Коровы",
			Answer3: "Собаки",
		}

		problem7 := Problem{
			VariantId: 3,
			Question: `Что в японском ресторане называется "осибори"?`,
			CorrectAnswer: "Влажное полотенце",
			Answer1: "Палочки для еды",
			Answer2: "Тарелка для суши",
			Answer3: "Соевый соус",
		}
		problem8 := Problem{
			VariantId: 3,
			Question: "Как называется человек, покоряющий крыши многоэтажных домов?",
			CorrectAnswer: "Руфер",
			Answer1: "Диггер",
			Answer2: "Сталкер",
			Answer3: "Байкер",
		}
		problem9 := Problem{
			VariantId: 3,
			Question: `К какому роду принадлежит панда 
			белого цвета с черными лапами, ушами 
			и участками вокруг глаз?`,
			CorrectAnswer: "Большая панда",
			Answer1: "Средняя панда",
			Answer2: "Малая панда",
			Answer3: "Белая панда",
		}


		tx, err := db.Sql.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockUser).Exec(
			user.Login,
			user.Password,
			user.LoginTime,
			user.LogoutTime,
		)
		if err != nil {
			return err
		}


		
		_, err = tx.Stmt(db.StmtMockVariant).Exec(
			variant1.Name,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockVariant).Exec(
			variant2.Name,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockVariant).Exec(
			variant3.Name,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem1.VariantId,
			problem1.Question,
			problem1.CorrectAnswer,
			problem1.Answer1,
			problem1.Answer2,
			problem1.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem2.VariantId,
			problem2.Question,
			problem2.CorrectAnswer,
			problem2.Answer1,
			problem2.Answer2,
			problem2.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem3.VariantId,
			problem3.Question,
			problem3.CorrectAnswer,
			problem3.Answer1,
			problem3.Answer2,
			problem3.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem4.VariantId,
			problem4.Question,
			problem4.CorrectAnswer,
			problem4.Answer1,
			problem4.Answer2,
			problem4.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem5.VariantId,
			problem5.Question,
			problem5.CorrectAnswer,
			problem5.Answer1,
			problem5.Answer2,
			problem5.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem6.VariantId,
			problem6.Question,
			problem6.CorrectAnswer,
			problem6.Answer1,
			problem6.Answer2,
			problem6.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem7.VariantId,
			problem7.Question,
			problem7.CorrectAnswer,
			problem7.Answer1,
			problem7.Answer2,
			problem7.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem8.VariantId,
			problem8.Question,
			problem8.CorrectAnswer,
			problem8.Answer1,
			problem8.Answer2,
			problem8.Answer3,
		)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(db.StmtMockProblem).Exec(
			problem9.VariantId,
			problem9.Question,
			problem9.CorrectAnswer,
			problem9.Answer1,
			problem9.Answer2,
			problem9.Answer3,
		)
		if err != nil {
			return err
		}


		tx.Commit()

	}
	
	return nil
}