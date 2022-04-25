--psql cselo -U cseloapp -f create-tables.sql



DROP TABLE IF EXISTS players CASCADE;

CREATE TABLE players (
 	id serial NOT NULL primary key,
 	initialname varchar(50) NOT NULL ,
 	steamid varchar(50) unique NOT NULL,
  profileid varchar(25),
  avatar varchar(200)
 	);


ALTER TABLE players
  OWNER TO cseloapp;


GRANT SELECT ON TABLE players TO cselojavaui;
GRANT UPDATE ON TABLE players TO cselojavaui;

DROP TABLE IF EXISTS matches CASCADE;

CREATE TABLE matches (
    id serial NOT NULL primary key,
    gamemode varchar(20),
    mapgroup varchar(40),
    mapfullname varchar(80),
    mapname varchar(40),
    scorea int default 0,
    scoreb int default 0,
    rounds int default 0,
    duration interval(6),
    matchstart timestamp,
    matchend timestamp,
    timestmp timestamp,
    completed boolean default false
);

ALTER TABLE matches
  OWNER TO cseloapp;

GRANT SELECT ON TABLE matches TO cselojavaui;


DROP TABLE IF EXISTS kills CASCADE;

CREATE TABLE kills (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    headshot boolean default false,
    weapon varchar(40),
    timestmp timestamp
);

ALTER TABLE kills
  OWNER TO cseloapp;

GRANT SELECT ON TABLE kills TO cselojavaui;


DROP TABLE IF EXISTS assists CASCADE;

CREATE TABLE assists (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    victim int REFERENCES players (id),
    timestmp timestamp
);

 ALTER TABLE assists
  OWNER TO cseloapp;

GRANT SELECT ON TABLE assists TO cselojavaui;


-- DROP TABLE IF EXISTS plantings CASCADE;

-- CREATE TABLE plantings (
--     id serial NOT NULL primary key,
--     match int REFERENCES matches (id),
--     actor int REFERENCES players (id),
--     timestmp timestamp
-- );

-- -- ALTER TABLE plantings
-- --  OWNER TO cseloapp;

-- GRANT SELECT ON TABLE plantings TO cselojavaui;


DROP TABLE IF EXISTS scoreaction CASCADE;

CREATE TABLE scoreaction (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    actiontype varchar(10),
    timestmp timestamp
);

ALTER TABLE scoreaction
  OWNER TO cseloapp;

GRANT SELECT ON TABLE scoreaction TO cselojavaui;


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

ALTER TABLE blindings
  OWNER TO cseloapp;

GRANT SELECT ON TABLE blindings TO cselojavaui;


DROP TABLE IF EXISTS grenadethrows CASCADE;

CREATE TABLE grenadethrows (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    grenadetype varchar(20),
    timestmp timestamp
);

ALTER TABLE grenadethrows
  OWNER TO cseloapp;

GRANT SELECT ON TABLE grenadethrows TO cselojavaui;


DROP TABLE IF EXISTS accolade CASCADE;

CREATE TABLE accolade (
    id serial NOT NULL primary key,
    match int REFERENCES matches (id),
    actor int REFERENCES players (id),
    accoladetype varchar(20),
    position int default 0,
    accoladevalue float,
    score float default 0,
    timestmp timestamp
);

ALTER TABLE accolade
  OWNER TO cseloapp;

GRANT SELECT ON TABLE accolade TO cselojavaui;


