-- TODO index all the tables

CREATE TYPE sex AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS users (
	id varchar(32) PRIMARY KEY,
	name varchar(64) NOT NULL,
	email varchar(320) UNIQUE,
	avatar_url text,

	google_id varchar(32) UNIQUE,

	active boolean NOT NULL DEFAULT true,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS parties (
	id int PRIMARY KEY,

	name varchar(64) NOT NULL,
	chairman varchar(64),
	established_date date,
	filing_date date,
	main_office_address text,
	mailing_address text,
	phone_number varchar(32),
	status varchar(32),

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS politicians (
	id SERIAL PRIMARY KEY,
	name varchar(64) NOT NULL,
	birthdate date,
	avatar_url text,
	sex sex,

	current_party_id int REFERENCES parties(id),

	meta jsonb,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS candidates (
	type varchar(32) NOT NULL,
	term int NOT NULL,
	politician_id int NOT NULL REFERENCES politicians(id),
	number int NOT NULL,
	elected boolean NOT NULL DEFAULT false,

	party_id int REFERENCES parties(id),

	area varchar(32),

	-- presidential candidates
	vice_president boolean NOT NULL DEFAULT false,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (type, term, politician_id)
);

CREATE TABLE IF NOT EXISTS legislators (
	politicians_id int NOT NULL REFERENCES politicians(id),
	term smallint NOT NULL,
	party_id int REFERENCES parties(id),
	onboard_date date,
	resign_date date,
	resign_reason text,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (politicians_id, term)
);

CREATE TABLE IF NOT EXISTS politician_questions (
	id SERIAL PRIMARY KEY,
	category varchar(32) NOT NULL,

	user_id varchar(32) NOT NULL REFERENCES users(id),
	question text NOT NULL,
	asked_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

	politician_id int NOT NULL REFERENCES politicians(id),
	reply text,
	replied_at timestamp,

	likes int NOT NULL DEFAULT 0,
	hidden boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS politician_question_likes (
	question_id int NOT NULL REFERENCES politician_questions(id),
	user_id varchar(32) NOT NULL REFERENCES users(id),
	PRIMARY KEY (question_id, user_id)
);

CREATE TABLE IF NOT EXISTS permissions (
	user_id varchar(32) NOT NULL REFERENCES users(id),
	resource varchar(32) NOT NULL,
	action varchar(128) NOT NULL,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (user_id, resource, action)
);

CREATE TABLE IF NOT EXISTS staging_data (
	id SERIAL PRIMARY KEY,
	table_name varchar(32) NOT NULL,
	fields jsonb NOT NULL,
	action varchar(32) NOT NULL,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
