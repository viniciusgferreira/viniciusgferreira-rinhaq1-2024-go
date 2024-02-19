-- Coloque scripts iniciais aqui
CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY,
    limite int,
    saldo int default 0,
    nome VARCHAR(100)
);

BEGIN;
INSERT INTO clientes (id, nome, limite)
VALUES
    (1, 'o barato sai caro', 1000 * 100),
    (2, 'zan corp ltda', 800 * 100),
    (3, 'les cruders', 10000 * 100),
    (4, 'padaria joia de cocaia', 100000 * 100),
    (5, 'kid mais', 5000 * 100);
COMMIT;

CREATE TABLE IF NOT EXISTS transacoes (
    id SERIAL PRIMARY KEY,
    cliente_id int,
    valor int,
    tipo VARCHAR(1),
    descricao VARCHAR(10),
    realizada_em TIMESTAMP,
    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
);

