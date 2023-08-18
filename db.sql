CREATE TABLE clicks (
    date date primary key,
    views bigint CHECK (views >= 0),
    clicks bigint CHECK (clicks >= 0),
    cost decimal CHECK (cost >= 0)
);

CREATE INDEX clicks_date_idx ON clicks (date);

grant all on clicks to statcounters;

select date::text, views, clicks, cost, round(cost / clicks, 2) as cpc, 
round(cost / views * 1000, 2) as cpm from clicks
where date between "2020-10-11" and "2020-10-11" order by date