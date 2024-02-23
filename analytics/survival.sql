WITH tb_daily AS (

    SELECT t1.IdCustomer,
        date(substr(t1.DtTransaction, 0,11)) AS DtTransaction,
        count(distinct t1.UUID) AS qtTransacao,
        sum(t1.Points) As qtPoints

    FROM tb_transactions AS t1

    LEFT JOIN tb_customers AS t2
    ON t1.idCustomer = t2.UUID
    
    WHERE t1.Points > 0
    AND t1.IdCustomer <> '5f8fcbe0-6014-43f8-8b83-38cf2f4887b3'
    AND t2.UUID is not null
    group by 1,2

),

tb_lag AS (

    SELECT *,
            lag(DtTransaction) OVER (PARTITION BY IdCustomer ORDER BY DtTransaction) AS lagDtTransaction
    FROM tb_daily

),

tb_user_day AS (

    SELECT IdCustomer,
            avg(julianday(DtTransaction) -julianday(lagDtTransaction)) AS qtDays
    FROM tb_lag
    -- where lagDtTransaction is not null
    group by 1
    
),

tb_soma_dias AS (

    SELECT qtDays,
        COUNT(IdCustomer) AS qt_user
    FROM tb_user_day
    GROUP BY 1

),

tb_survival AS (

    SELECT *,
            1.0 * qt_user / (SELECT SUM(qt_user) FROM tb_soma_dias) as pct,
            SUM(1.0 * qt_user / (SELECT SUM(qt_user) FROM tb_soma_dias)) OVER (PARTITION BY 1 ORDER BY qtDays) as pct_acum,
            sum(qt_user)  OVER (PARTITION BY 1 ORDER BY qtDays) as qtd_acum

    FROM tb_soma_dias

),

tb_qtde_dias AS (

    SELECT t1.idCustomer,
           count(DISTINCT date(substr(t1.DtTransaction, 0,11))) as qtde_days
    FROM tb_transactions AS t1

    LEFT JOIN tb_customers AS t2
    ON t1.idCustomer = t2.UUID

    WHERE t1.Points > 0
    AND t1.IdCustomer <> '5f8fcbe0-6014-43f8-8b83-38cf2f4887b3'
    AND t2.UUID is not null
    group by 1

)

-- select sum(qt_user) from tb_soma_dias -- 273

-- select sum(case when qtde_days > 1 then 1 else 0 end),
--        avg(case when qtde_days > 1 then 1 else 0 end)

-- from tb_qtde_dias

SELECT *
FROM tb_lag