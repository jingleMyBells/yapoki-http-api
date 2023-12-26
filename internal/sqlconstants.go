package internal

const (
	schemaVariantSQL = `
	CREATE TABLE IF NOT EXISTS variant (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL
	);
	`
	schemaProblemSQL = `
	CREATE TABLE IF NOT EXISTS problem (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		variant_id INT NOT NULL,
		question VARCHAR(1000) NOT NULL,
		correct_answer VARCHAR(1000) NOT NULL,
		answer_1 VARCHAR(1000) NOT NULL,
		answer_2 VARCHAR(1000) NOT NULL,
		answer_3 VARCHAR(1000) NOT NULL,
		FOREIGN KEY (variant_id)  REFERENCES variant (id)
	);
	`
	schematTestSQL = `
	CREATE TABLE IF NOT EXISTS test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL,
		variant_id INT NOT NULL,
		start_time TIMESTAMP,
		FOREIGN KEY (user_id)  REFERENCES user (id),
		FOREIGN KEY (variant_id)  REFERENCES variant (id)
	);
	`

	schematTestAnswersSQL = `
	CREATE TABLE IF NOT EXISTS test_answer (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		problem_id INT NOT NULL,
		test_id INT NOT NULL,
		answer VARCHAR(1000) NOT NULL,
		FOREIGN KEY (test_id)  REFERENCES test (id),
		FOREIGN KEY (problem_id)  REFERENCES problem (id)
	);
	`

	schemaTestResultsSQL = `
	CREATE TABLE IF NOT EXISTS test_result (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_id INT NOT NULL,
		percent INT NOT NULL,
		FOREIGN KEY (test_id)  REFERENCES test (id)
	);
	`

	schemaUserSQL = `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		login_time TIMESTAMP,
		logout_time TIMESTAMP,
		cookie VARCHAR(255) UNIQUE
	);
	`

	insertNewTestSQL = `
	INSERT INTO test (
		user_id, variant_id, start_time
	) VALUES (
		?, ?, ?
	)
	`

	insertNewTestAnswers = `
	INSERT INTO test_answer (
		test_id, problem_id, answer
	) VALUES (
		?, ?, ?
	)
	`

	insertNewTestResult = `
	INSERT INTO test_result (
		test_id, percent
	) VALUES (
		?, ?
	)
	`

	insertMockUser = `
	INSERT INTO user (
		login, password, login_time, logout_time
	) VALUES (
		?, ?, ?, ?
	)
	`
)


// список доступных вариантов: variant (id, name)
// список заданий для варианта: problem (id, variant_id, question, correct_answer, answers)
// тестирование: test (id, user_id, variant_id, timestamp)
// список ответов пользователя: test_answers (id, test_id, problem_id, answer)
// результаты: test_results (id, test_id, percent)
// user: user (id, login, password, login_time, logout_time)


// schemaSQL = `
// CREATE TABLE IF NOT EXISTS tests (
// 	time TIMESTAMP,
// 	symbol VARCHAR(32),
// 	price FLOAT,
// 	buy BOOLEAN
// );


// insertSQL = `
// INSERT INTO tests (
// 	time, symbol, price, buy
// ) VALUES (
// 	?, ?, ?, ?
// )


// db, err := sql.Open("postgres", "user=username password=password host=localhost dbname=mydb sslmode=disable") 
// if err != nil { 
// 	log.Fatalf("Error: Unable to connect to database: %v", err) 
// 	} 
// defer db.Close()

// rows, err := db.Query("SELECT id, name FROM users") 
// if err != nil {
// 	 log.Fatalf("Error: Unable to execute query: %v", err) 
// } 
// defer rows.Close() 
// for rows.Next() { 
// 	var id int64 
// 	var name string 
// 	rows.Scan(&id, &name) 
// 	fmt.Printf("User ID: %d, Name: %s\n", id, name) } 


// var id int64 
// var name string 
// row := db.QueryRow("SELECT id, name FROM users WHERE id = $1", userID) 
// if err := row.Scan(&id, &name); err == sql.ErrNoRows {
// 	 fmt.Println("Пользователь не найден") 
// } else if err != nil {
// 	 log.Fatalf("Error: Unable to execute query: %v", err) 
// } 
//   else { 
// 	fmt.Printf("User ID: %d, Name: %s\n", id, name) 
// }


// result, err := db.Exec("UPDATE users SET email = $1 WHERE id = $2", email, userID) 
// if err != nil {
// 	 log.Fatalf("Error: Unable to execute update: %v", err) 
// } 
// affectedRows, _ := result.RowsAffected() 
// fmt.Printf("Updated %d rows\n", affectedRows)