DROP TABLE IF EXISTS tb_transactions_cart;
CREATE TABLE IF NOT EXISTS tb_transactions_cart (
 UUID VARCHAR(30),
 IdTransaction VARCHAR(100),
 Product VARCHAR(250),
 Quantity INT
);