package internal

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)




type DB struct {
	Sql *sql.DB
	StmtTest *sql.Stmt
	StmtTestAnswer *sql.Stmt
	StmtTestResult *sql.Stmt
	StmtMockUser *sql.Stmt
	StmtMockVariant *sql.Stmt
	StmtMockProblem *sql.Stmt
}

var db *DB

func init() {
	database, err := NewDB("db.sqlite3")
	if err != nil {
		log.Printf("Ошибка подключения к базе: %v", err)
	}
	db = database
}

func GetDB() *DB {
	return db
}

func NewDB(dbFile string) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	"localhost", 5432, "postgres", "12345", "postgres")
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaUserSQL); err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaVariantSQL); err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaProblemSQL); err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schematTestSQL); err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schematTestAnswersSQL); err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaTestResultsSQL); err != nil {
		return nil, err
	}

	stmtTest, err := sqlDB.Prepare(insertNewTestSQL)
	if err != nil {
		return nil, err
	}

	stmtTestAnswer, err := sqlDB.Prepare(insertNewTestAnswers)
	if err != nil {
		return nil, err
	}

	stmtTestResult, err := sqlDB.Prepare(insertNewTestResult)
	if err != nil {
		return nil, err
	}

	// ниже стейтменты для тестовые данных в бд
	stmtMockUser, err := sqlDB.Prepare(insertMockUser)
	if err != nil {
		return nil, err
	}

	stmtMockVariant, err := sqlDB.Prepare(insertMockVariant)
	if err != nil {
		return nil, err
	}

	stmtMockProblem, err := sqlDB.Prepare(insertMockProblem)
	if err != nil {
		return nil, err
	}

	db := DB {
		Sql: sqlDB,
		StmtTest: stmtTest,
		StmtTestAnswer: stmtTestAnswer,
		StmtTestResult: stmtTestResult,
		StmtMockUser: stmtMockUser,
		StmtMockVariant: stmtMockVariant,
		StmtMockProblem: stmtMockProblem,
	}

	return &db, nil
}


func (db *DB) AddTest(test Test) error {
	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.StmtTest).Exec(
		test.UserId,
		test.VariantId,
		test.StartTime,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}


func (db *DB) AddTestAnswer(testAnswer TestAnswer) error {
	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.StmtTestAnswer).Exec(
		testAnswer.TestId,
		testAnswer.ProblemId,
		testAnswer.Answer,
	)
	
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (db *DB) AddTestResult(testResult TestResult) error {
	tx, err := db.Sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.StmtTestResult).Exec(
		testResult.TestId,
		testResult.Percent,
	)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}


func (db *DB) Close() error {
	defer func() {
		db.StmtTest.Close()
		db.StmtTestAnswer.Close()
		db.StmtTestResult.Close()
		db.StmtMockUser.Close()
		db.StmtMockVariant.Close()
		db.StmtMockProblem.Close()
		db.Sql.Close()
	}()

	return nil
}
