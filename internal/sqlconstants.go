package internal

const (
	schemaVariantSQL = `
	CREATE TABLE IF NOT EXISTS variant (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	`
	schemaProblemSQL = `
	CREATE TABLE IF NOT EXISTS problem (
		id SERIAL PRIMARY KEY,
		variant_id INT NOT NULL,
		question VARCHAR(1000) NOT NULL,
		correct_answer VARCHAR(1000) NOT NULL,
		answer_1 VARCHAR(1000) NOT NULL,
		answer_2 VARCHAR(1000) NOT NULL,
		answer_3 VARCHAR(1000) NOT NULL,
		FOREIGN KEY (variant_id)  REFERENCES variant (id) ON DELETE CASCADE
	);
	`
	schematTestSQL = `
	CREATE TABLE IF NOT EXISTS test (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		variant_id INT NOT NULL,
		start_time TIMESTAMP,
		FOREIGN KEY (user_id)  REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (variant_id)  REFERENCES variant (id) ON DELETE CASCADE
	);
	`

	schematTestAnswersSQL = `
	CREATE TABLE IF NOT EXISTS test_answer (
		id SERIAL PRIMARY KEY,
		problem_id INT NOT NULL,
		test_id INT NOT NULL,
		answer VARCHAR(1000) NOT NULL,
		FOREIGN KEY (test_id)  REFERENCES test (id) ON DELETE CASCADE,
		FOREIGN KEY (problem_id)  REFERENCES problem (id) ON DELETE CASCADE
	);
	`

	schemaTestResultsSQL = `
	CREATE TABLE IF NOT EXISTS test_result (
		id SERIAL PRIMARY KEY,
		test_id INT NOT NULL,
		percent INT NOT NULL,
		FOREIGN KEY (test_id)  REFERENCES test (id) ON DELETE CASCADE
	);
	`

	schemaUserSQL = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
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
		$1, $2, $3
	)
	`

	insertNewTestAnswers = `
	INSERT INTO test_answer (
		test_id, problem_id, answer
	) VALUES (
		$1, $2, $3
	)
	`

	insertNewTestResult = `
	INSERT INTO test_result (
		test_id, percent
	) VALUES (
		$1, $2
	)
	`

	insertMockUser = `
	INSERT INTO users (
		login, password, login_time, logout_time
	) VALUES (
		$1, $2, $3, $4
	)
	`

	insertMockVariant = `
	INSERT INTO variant (
		name
	) VALUES (
		$1
	)
	`

	insertMockProblem = `
	INSERT INTO problem (
		variant_id, question, correct_answer, answer_1, answer_2, answer_3
	) VALUES (
		$1, $2, $3, $4, $5, $6
	)
	`
)
