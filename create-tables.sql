--psql cselo -U cseloapp -f create-tables.sql

DROP TABLE IF EXISTS players CASCADE;

CREATE TABLE players (
 	id serial NOT NULL primary key,
 	initialname varchar(50) NOT NULL ,
 	steamid varchar(50) unique NOT NULL
 	);

--ALTER TABLE players
--  OWNER TO cselo;


DROP TABLE IF EXISTS kills CASCADE;

CREATE TABLE kills (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    headshot boolean,
    weapon varchar(40),
    timestmp timestamp
);

-- ALTER TABLE kills
--  OWNER TO cseloapp;

DROP TABLE IF EXISTS assists CASCADE;

CREATE TABLE assists (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    timestmp timestamp
);

-- ALTER TABLE assists
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS plantings CASCADE;

CREATE TABLE plantings (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    timestmp timestamp
);

-- ALTER TABLE plantings
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS rescues CASCADE;

CREATE TABLE rescues (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    timestmp timestamp
);

-- ALTER TABLE rescues
--  OWNER TO cseloapp;


-- ALTER TABLE assists
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS bombings CASCADE;

CREATE TABLE bombings (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    timestmp timestamp
);

-- ALTER TABLE planbombingstings
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS defuses CASCADE;

CREATE TABLE defuses (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    timestmp timestamp
);

-- ALTER TABLE defuses
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS blindings CASCADE;

CREATE TABLE blindings (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    duration float,
    victimtype varchar(10),
    timestmp timestamp
);

-- ALTER TABLE blindings
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS grenadethrows CASCADE;

CREATE TABLE grenadethrows (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    grenadetype varchar(20),
    timestmp timestamp
);

-- ALTER TABLE grenadethrows
--  OWNER TO cseloapp;


DROP TABLE IF EXISTS matches CASCADE;

CREATE TABLE matches (
    id serial NOT NULL primary key,
    gamemode varchar(20),
    mapgroup varchar(40),
    mapfullname varchar(80),
    mapname varchar(40),
    scorea int,
    scoreb int,
    duration interval(6),
    matchstart timestamp,
    matchend timestamp,
    timestmp timestamp
);

-- ALTER TABLE matches
--  OWNER TO cseloapp;


CREATE TABLE rounds (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    duration interval(6),
    roundstart timestamp,
    roundend timestamp,
    timestmp timestamp
);

-- ALTER TABLE rounds
--  OWNER TO cseloapp;
