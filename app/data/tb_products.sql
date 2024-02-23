DROP TABLE IF EXISTS tb_products;
CREATE TABLE IF NOT EXISTS tb_products (
    idProduct int,
    descProducts varchar(150),
    vlProduct int,
    PtProduct int
);

INSERT INTO tb_products VALUES (1, 'Pontos Gerais', 0, 1);
INSERT INTO tb_products VALUES (2, 'Resgate de Pontos', 0, -1);