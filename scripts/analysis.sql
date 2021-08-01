\! echo "Kills:";
select initialname,count(actor) as kills,count(case when headshot=true then 1 end) as headshot, round(cast(count(case when headshot=true then 1 end) as float)/count(actor) * 1000)/10 as "hs%"  from kills
left join players on actor=players.id
WHERE timestmp > current_date - interval '3' day
group by initialname
order by count(actor) DESC;

\! echo "All grenades:";
SELECT initialname,count(case when grenadetype='flashbang' then 1 end) as flash,
    count(case when grenadetype='hegrenade' then 1 end) as he,
	count(case when grenadetype='molotov' then 1 end) as molotov,
	count(case when grenadetype='smokegrenade' then 1 end) as smoke,
	count(case when grenadetype='decoy' then 1 end) as decoy FROM grenadethrows
LEFT JOIN players ON actor=players.id
WHERE timestmp > current_date - interval '3' day
GROUP BY initialname
ORDER BY flash DESC;


\! echo "Flashes:";
select initialname,count(case when victimtype='enemy' then 1 end) as enemyflashes, count(case when victimtype='teammate' then 1 end) as teammateflash, count(case when victimtype='self' then 1 end) as selfflash   from blindings
left join players on actor=players.id
WHERE timestmp > current_date - interval '3' day
group by initialname
order by count(case when victimtype='enemy' then 1 end) DESC;


\! echo "Special actions:";
SELECT initialname,count(case when actiontype='planting' then 1 end) as planting,
    count(case when actiontype='bombing' then 1 end) as bombing,
	count(case when actiontype='defuse' then 1 end) as defuse,
	count(case when actiontype='rescue' then 1 end) as rescue FROM scoreaction
LEFT JOIN players ON actor=players.id
WHERE timestmp > current_date - interval '3' day
GROUP BY initialname
ORDER BY planting DESC;
