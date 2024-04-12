CREATE SCHEMA IF NOT EXISTS accounting;
GRANT ALL ON SCHEMA accounting TO postgres;

CREATE TABLE IF NOT EXISTS accounting.currencies(
	id   SERIAL       PRIMARY KEY,
	iso  VARCHAR(3)   NOT NULL UNIQUE,
	name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounting.accounts(
	id          SERIAL PRIMARY KEY,
	description TEXT   NOT NULL
);

CREATE TABLE IF NOT EXISTS accounting.documents(
	id           SERIAL       PRIMARY KEY,
	date         VARCHAR(255) NOT NULL,
	posting_date VARCHAR(255) NOT NULL,
	reference    VARCHAR(255) NOT NULL,
	description  TEXT         NOT NULL,
	currency_id  INTEGER      NOT NULL REFERENCES accounting.currencies(id)
);

CREATE TABLE IF NOT EXISTS accounting.document_position_types(
	id          SERIAL  PRIMARY KEY,
	document_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS accounting.document_positions(
	id          SERIAL  PRIMARY KEY,
	document_id SERIAL  NOT NULL REFERENCES accounting.documents(id),
	account_id  SERIAL  NOT NULL REFERENCES accounting.accounts(id),
	description TEXT    NOT NULL,
	type_id     INTEGER NOT NULL,
	amount      NUMERIC NOT NULL
);
