# Artifacts

Small application that counts some data in MySQL database (this application call it artifacts) and then send it 
to InfluxDB using UDP.

## Preparation

1. Download and compile application
2. Create database table holding requests to made:
```sql
CREATE TABLE `__artifacts` (
  `name` varchar(1024) CHARACTER SET latin1 NOT NULL,
  `query` varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL,
  `interval` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
```
3. Fill some data to table

## Query format

Application works with almost any query type that meets following criteria
- Metric value must be integer
- Metric value must be last column
- One metric value per query
- Grouping columns must be first

Valid queries:
```sql
SELECT count(1) FROM foo;
SELECT count(1) FROM foo WHERE bar='baz';
SELECT g1, g2, count(1) FROM foo GROUP BY g1, g2;
```

Not supported queries:
```sql
SELECT min(f), max(f) FROM foo;           -- Two values in single query
SELECT avg(f) FROM foo;                   -- Possibly non integer value
SELECT count(1), g1 FROM foo GROUP BY g1; -- Value column is not last
```

## Testing

To test that everything is ok, run command `app test`.
This will establish connection to MySQL server and performs all queries in sequential order.
No data will be sent to InfluxDB.

For further information and configuration options run `app test --help`

## Running daemon

`app serve` to run server, `app serve --help` for more information and configuration options.