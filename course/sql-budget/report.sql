.headers on
.mode column

WITH monthly AS (
    SELECT c.kind, SUM(t.amount_tiyn) AS amount_tiyn
    FROM transactions AS t
    JOIN categories AS c ON c.id = t.category_id
    WHERE t.happened_on >= '2026-09-01'
      AND t.happened_on <  '2026-10-01'
    GROUP BY c.kind
)
SELECT
    printf('%.2f', COALESCE(SUM(amount_tiyn) FILTER (WHERE kind = 'income'), 0) / 100.0) AS income_kzt,
    printf('%.2f', COALESCE(SUM(amount_tiyn) FILTER (WHERE kind = 'expense'), 0) / 100.0) AS expense_kzt,
    printf('%.2f', (
        COALESCE(SUM(amount_tiyn) FILTER (WHERE kind = 'income'), 0) -
        COALESCE(SUM(amount_tiyn) FILTER (WHERE kind = 'expense'), 0)
    ) / 100.0) AS saved_kzt
FROM monthly;

SELECT c.name,
       printf('%.2f', SUM(t.amount_tiyn) / 100.0) AS spent_kzt,
       printf('%.2f', c.monthly_limit_tiyn / 100.0) AS limit_kzt,
       CASE WHEN SUM(t.amount_tiyn) > c.monthly_limit_tiyn THEN 'OVER' ELSE 'OK' END AS status
FROM transactions AS t
JOIN categories AS c ON c.id = t.category_id
WHERE c.kind = 'expense'
  AND t.happened_on >= '2026-09-01' AND t.happened_on < '2026-10-01'
GROUP BY c.id, c.name, c.monthly_limit_tiyn
ORDER BY SUM(t.amount_tiyn) DESC;

