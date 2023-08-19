CREATE TABLE clicks (
    date date primary key,
    views bigint CHECK (views >= 0),
    clicks bigint CHECK (clicks >= 0),
    cost numeric CHECK (cost >= 0) CHECK (scale(cost) <= 2)
);

CREATE INDEX clicks_date_idx ON clicks (date);