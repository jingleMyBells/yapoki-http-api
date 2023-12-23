package internal

import (
	"database/sql"
	// "errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)




type DB struct {
	Sql *sql.DB
	StmtTest *sql.Stmt
	StmtTestAnswer *sql.Stmt
	StmtTestResult *sql.Stmt
	Buffer []Test
}

var db *DB

func init() {
	database, err := NewDB("db.sqlite3")
	if err != nil {
		fmt.Println(err)
	}
	db = database
}

func GetDB() *DB {
	return db
}

func NewDB(dbFile string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
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

	if _, err = sqlDB.Exec(schemaUserSQL); err != nil {
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

	db := DB {
		Sql: sqlDB,
		StmtTest: stmtTest,
		StmtTestAnswer: stmtTestAnswer,
		StmtTestResult: stmtTestResult,
		Buffer: make([]Test, 0, 1),
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
		db.Sql.Close()
	}()

	return nil
}


// func (db *DB) Add(test Test) error {
// 	if len(db.Buffer) == cap(db.Buffer) {
// 		return errors.New("tests buffer is full")
// 	}

// 	db.Buffer = append(db.Buffer, test)
// 	if len(db.Buffer) == cap(db.Buffer) {
// 		if err := db.Flush(); err != nil {
// 			return fmt.Errorf("не получилось слить в базу: %w", err)
// 		}
// 	}

// 	return nil
// }



// func (db *DB) Flush() error {
// 	tx, err := db.Sql.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	for _, test := range db.Buffer {
// 		_, err = tx.Stmt(db.Stmt).Exec(test.Time, test.Symbol, test.Price, test.isBuy)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	db.Buffer = db.Buffer[:0]
// 	return tx.Commit()
// }

// func (db *DB) Close() error {
// 	defer func() {
// 		db.Stmt.Close()
// 		db.Sql.Close()
// 	}()

// 	if err := db.Flush(); err != nil {
// 		return err
// 	}

// 	return nil
// }



// func main() {
// 	fmt.Println("Что-то вообще происходит?")

// 	db, err := NewDB("db.sqlite3")
// 	if err != nil {
// 		fmt.Println("Какая-то ошибкаЖ", err)
// 	}

// 	test := Test{
// 		Time: time.Now(),
// 		Symbol: "lalala",
// 		Price: 3.14,
// 		isBuy: true,
// 	}

// 	err = db.Add(test)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	err = db.Add(test)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	err = db.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Выход")
	
// }