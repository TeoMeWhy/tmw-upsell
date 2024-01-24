DROP TABLE IF EXISTS tb_customers;
CREATE TABLE IF NOT EXISTS tb_customers (
 UUID varchar(30),
 Name varchar(100),
 Email varchar(100),
 CPF varchar(11),
 Points int
);