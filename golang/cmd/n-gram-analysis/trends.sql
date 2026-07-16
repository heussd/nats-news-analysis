-- Upcoming / rising trend words
-- Compares each n-gram's frequency in a RECENT window against a preceding
-- BASELINE window and ranks phrases by their growth. Brand-new phrases
-- (no baseline) surface first.
WITH params AS (
    SELECT
        (SELECT max(timestamp) FROM public.ngrams) AS ref_ts,
        INTERVAL '7 days'  AS recent_window,
        INTERVAL '30 days' AS baseline_window
),
recent AS (
    SELECT words, n_gram, SUM(count) AS recent_count
    FROM public.ngrams, params
    WHERE timestamp > ref_ts - recent_window
    GROUP BY words, n_gram
),
baseline AS (
    SELECT words, n_gram, SUM(count) AS baseline_count
    FROM public.ngrams, params
    WHERE timestamp <= ref_ts - recent_window
      AND timestamp >  ref_ts - recent_window - baseline_window
    GROUP BY words, n_gram
)
SELECT
    r.words,
    r.n_gram,
    r.recent_count,
    COALESCE(b.baseline_count, 0)                     AS baseline_count,
    r.recent_count - COALESCE(b.baseline_count, 0)    AS absolute_growth,
    ROUND(r.recent_count::numeric
          / NULLIF(b.baseline_count, 0), 2)           AS growth_ratio
FROM recent r
LEFT JOIN baseline b USING (words, n_gram)
WHERE r.n_gram >= 2            -- phrases; change to r.n_gram = 1 for single words
  AND r.recent_count >= 2 -- ignore one-off noise
ORDER BY
    (b.baseline_count IS NULL) DESC, -- brand-new phrases first
    growth_ratio DESC NULLS FIRST,
    absolute_growth DESC
LIMIT 50;