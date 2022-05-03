package datastore

const CREATE_USERS_TABLE = `CREATE TABLE IF NOT EXISTS USERS(
	ID int primary key auto_increment,
	USER_NAME text,
	ACCOUNT text,
	AMOUNT int default 0
)`

const CREATE_ORDERS_TABLE = `CREATE TABLE IF NOT EXISTS ORDERS(
	ID int primary key auto_increment,
	ACCOUNT text,
	PRODUCT_NAME text,
	PRICE int,
	ORDER_STATUS text
)`
