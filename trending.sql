-- Upcoming / rising trend words
-- Compares each n-gram's frequency in a RECENT window against a preceding
-- BASELINE window and ranks phrases by their growth. Frequency is the number
-- of occurrences (rows).
--
-- How it works:
--   * recent / baseline CTEs count occurrences (COUNT(*)) of each phrase per
--     window, anchored to the latest timestamp in the table (ref_ts).
--   * Counts are normalised to a PER-DAY RATE so the differently sized windows
--     (7 vs 30 days) are directly comparable.
--   * growth_ratio uses Laplace smoothing so tiny baselines don't produce
--     absurdly large ratios.
--   * absolute_growth = raw increase in occurrences between the windows.
--   * z_score approximates statistical significance: how many standard
--     deviations the recent count is above the count expected from the
--     baseline rate (Poisson assumption). Higher = less likely to be noise.
--   * Brand-new phrases (no baseline) are surfaced first and require a stricter
--     minimum count to keep one-off noise out.
--
-- Knobs to tune:
--   * The '7 days' / '30 days' recent/baseline windows (keep *_days in sync).
--   * The n_gram filter (set = 1 for single words).
--   * min_recent / min_recent_new noise thresholds.
--   * Uncomment the language / source filters for per-language/source trends.
WITH params AS
    (SELECT
         (SELECT max("timestamp")
          FROM public.ngrams) AS ref_ts,
            INTERVAL '7 days' AS recent_window,
            INTERVAL '30 days' AS baseline_window,
            7.0::numeric AS recent_days, -- keep in sync with recent_window
 30.0::numeric AS baseline_days, -- keep in sync with baseline_window
 2 AS min_recent, -- min occurrences for known phrases
 3 AS min_recent_new -- min occurrences for brand-new phrases
),
     recent AS
    (SELECT words,
            n_gram,
            COUNT(*) AS recent_count,
            COUNT(DISTINCT source) AS recent_sources
     FROM public.ngrams,
          params
     WHERE "timestamp" > ref_ts - recent_window -- AND language = 'en'
 -- AND source = '...'

     GROUP BY words,
              n_gram),
     baseline AS
    (SELECT words,
            n_gram,
            COUNT(*) AS baseline_count
     FROM public.ngrams,
          params
     WHERE "timestamp" <= ref_ts - recent_window
         AND "timestamp" > ref_ts - recent_window - baseline_window -- AND language = 'en'
 -- AND source = '...'

     GROUP BY words,
              n_gram),
     joined AS
    (SELECT r.words,
            r.n_gram,
            r.recent_count,
            r.recent_sources,
            COALESCE(b.baseline_count, 0) AS baseline_count,
            (b.baseline_count IS NULL) AS is_new, -- per-day rates make the two window sizes comparable
 r.recent_count::numeric / p.recent_days AS recent_rate,
 COALESCE(b.baseline_count, 0)::numeric / p.baseline_days AS baseline_rate, -- expected recent count if the baseline rate simply continued
 (COALESCE(b.baseline_count, 0)::numeric / p.baseline_days) * p.recent_days AS expected_recent
     FROM recent r
     LEFT JOIN baseline b USING (words,
                                 n_gram)
     CROSS JOIN params p
     WHERE r.n_gram >= 2 -- phrases; change to = 1 for single words

         AND r.recent_count >= CASE
                                   WHEN b.baseline_count IS NULL THEN p.min_recent_new
                                   ELSE p.min_recent
                               END)
SELECT words,
       n_gram,
       recent_count,
       recent_sources,
       baseline_count,
       ROUND(recent_rate, 3) AS recent_per_day,
       ROUND(baseline_rate, 3) AS baseline_per_day,
       recent_count - baseline_count AS absolute_growth, -- Laplace-smoothed per-day rate ratio: robust against tiny baselines
 ROUND((recent_rate + (1.0 / 7)) / (baseline_rate + (1.0 / 30)), 2) AS growth_ratio, -- Poisson z-score vs. the count expected from the baseline rate
 ROUND((recent_count - expected_recent) / SQRT(GREATEST(expected_recent, 1.0)), 2) AS z_score
FROM joined
ORDER BY is_new DESC, -- brand-new phrases first
 z_score DESC, -- then most statistically surprising
 recent_sources DESC, -- prefer phrases corroborated across many distinct sources
 growth_ratio DESC,
 absolute_growth DESC
LIMIT 500;