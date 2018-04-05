CREATE TABLE users (
    id          SERIAL PRIMARY KEY ,
    name        VARCHAR(20)
);

CREATE TABLE relationships (
    id  SERIAL PRIMARY KEY ,
    suid bigint,
    tuid bigint,
    status varchar(10)
    UNIQUE (suid, tuid)
)