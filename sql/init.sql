CREATE TABLE users (
    id          SERIAL PRIMARY KEY ,
    likeset     BIGINT[],
    dislikeset  BIGINT[],
    match       BIGINT[],
    name        VARCHAR(20)
);
