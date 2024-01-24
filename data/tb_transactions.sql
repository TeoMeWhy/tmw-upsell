DROP TABLE IF EXISTS tb_transactions;
CREATE TABLE IF NOT EXISTS tb_transactions (
 UUID VARCHAR(30),
 IdCustomer VARCHAR(100),
 DtTransaction datetime,
 Points INT
);
