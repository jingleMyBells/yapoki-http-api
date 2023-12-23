package internal


// schemaSQL = `
// CREATE TABLE IF NOT EXISTS tests (
// 	time TIMESTAMP,
// 	symbol VARCHAR(32),
// 	price FLOAT,
// 	buy BOOLEAN
// );

const (
	schemaVariantSQL = `
	CREATE TABLE IF NOT EXISTS variant (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255)
	);
	`
	schemaProblemSQL = `
	CREATE TABLE IF NOT EXISTS problem (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		variant_id INT,
		question VARCHAR(1000),
		correct_answer VARCHAR(1000),
		answer_1 VARCHAR(1000),
		answer_2 VARCHAR(1000),
		answer_3 VARCHAR(1000),
		FOREIGN KEY (variant_id)  REFERENCES variant (id)
	);
	`
	schematTestSQL = `
	CREATE TABLE IF NOT EXISTS test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT,
		variant_id INT,
		start_time TIMESTAMP,
		FOREIGN KEY (user_id)  REFERENCES user (id),
		FOREIGN KEY (variant_id)  REFERENCES variant (id)
	);
	`

	schematTestAnswersSQL = `
	CREATE TABLE IF NOT EXISTS test_answer (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		problem_id INT,
		test_id INT,
		answer VARCHAR(1000),
		FOREIGN KEY (test_id)  REFERENCES test (id),
		FOREIGN KEY (problem_id)  REFERENCES problem (id)
	);
	`

	schemaTestResultsSQL = `
	CREATE TABLE IF NOT EXISTS test_result (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_id INT,
		percent INT,
		FOREIGN KEY (test_id)  REFERENCES test (id)
	);
	`

	schemaUserSQL = `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login VARCHAR(255),
		password VARCHAR(255),
		login_time TIMESTAMP,
		logout_time TIMESTAMP
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