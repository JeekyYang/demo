CREATE TABLE "user" (
    id          BIGINT NOT NULL PRIMARY KEY ,
    likeset     BIGINT[],
    dislikeset  BIGINT[],
    match       BIGINT[],
    name        VARCHAR(20)
);
