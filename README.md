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

## Testing

To test that everything is ok, run command `app test`.
This will establish connection to MySQL server and performs all queries in sequential order.
No data will be sent to InfluxDB.

For further information and configuration options run `app test --help`

## Running daemon

`app serve` to run server, `app serve --help` for more information and configuration options.