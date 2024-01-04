WITH RECURSIVE Counter AS (
    SELECT 1 AS Value
    UNION ALL
    SELECT Value + 1 FROM Counter WHERE Value < 5
),
DummyData AS (
    SELECT 
        CASE WHEN Counter.Value % 2 = 0 THEN 'income' ELSE 'expense' END AS type,
        'Description ' || Counter.Value AS description,
        ABS(RANDOM()) % (1000 - 50) + 10 AS amount
    FROM Counter
)
INSERT INTO transactions (type, description, amount)
SELECT * FROM DummyData;
