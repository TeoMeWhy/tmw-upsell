DROP TABLE IF EXISTS tb_customers;
CREATE TABLE IF NOT EXISTS tb_customers (
    UUID varchar(30),
    IdOrg VARCHAR(100),
    Name varchar(100),
    Email varchar(100),
    CPF varchar(11),
    Points int,
    TelResidencial varchar(20),
    TelComercial varchar(20),
    Instagram varchar(150)
);

CREATE INDEX ix_UUID
ON tb_customers (UUID);