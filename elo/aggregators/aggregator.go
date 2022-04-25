package aggregators

type Aggregator interface {
}

/*
SELECT initialname, date_trunc('day', kills.timestmp) AS day,
   ROUND(CAST(COUNT(CASE WHEN actor=players.id THEN 1 END) AS NUMERIC) / GREATEST(COUNT(distinct match),1),2) AS killspermatch,
   ROUND(CAST(COUNT(CASE WHEN actor=players.id THEN 1 END) AS NUMERIC) / GREATEST(COUNT(CASE WHEN victim=players.id THEN 1 END),1),2) AS kdratio
FROM kills LEFT JOIN matches on (kills.match=matches.id)
LEFT JOIN players ON (kills.actor=players.id OR kills.victim=players.id)
WHERE matches.completed=true AND players.id=2
GROUP BY day, initialname;
*/
