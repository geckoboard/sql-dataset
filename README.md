# SQL-Dataset, by Geckoboard

[![CircleCI](https://circleci.com/gh/geckoboard/sql-dataset.svg?style=svg)](https://circleci.com/gh/geckoboard/sql-dataset)  [![codecov](https://codecov.io/gh/geckoboard/sql-dataset/branch/master/graph/badge.svg)](https://codecov.io/gh/geckoboard/sql-dataset)

Quickly and easily send data from Microsoft SQL Server, MySQL, Postgres and SQLite databases to Geckoboard Datasets.

SQL-Dataset is a command line app that takes the hassle out of integrating your database with Geckoboard. Rather than having to work with client libraries and write a bunch of code to connect to and query your database, with SQL-Dataset all you need to do is fill out a simple config file.

SQL-Dataset is available for macOS, Linux, and Windows.

## Quickstart

### 1. Download the app

* [macOS](https://github.com/geckoboard/sql-dataset/releases/download/v0.2.4/sql-dataset-darwin-10.10-amd64)
* Linux [x86](https://github.com/geckoboard/sql-dataset/releases/download/v0.2.4/sql-dataset-linux-386) / [x64](https://github.com/geckoboard/sql-dataset/releases/download/v0.2.4/sql-dataset-linux-amd64)
* Windows [x86](https://github.com/geckoboard/sql-dataset/releases/download/v0.2.4/sql-dataset-windows-8.0-386.exe) / [x64](https://github.com/geckoboard/sql-dataset/releases/download/v0.2.4/sql-dataset-windows-8.0-amd64.exe)

#### Make it executable (macOS / Linux)

On macOS and Linux you'll need to open a terminal and run `chmod u+x path/to/file` (replacing `path/to/file` with the actual path to your downloaded app) in order to make the app executable.

### 2. Create a config file

SQL-Datasets works by reading all of the information it needs from a YAML file. We've prepared an [example one](example.yml) for you so you can get started quickly. The fields are fairly self-explanatory, but you can learn more about them [below](README.md#building-your-config-file).

### 3. Run the script

Make sure that the SQL-Dataset app and your config file are in the same folder, then from the command line navigate to that folder and run

```
./sql-dataset -config config.yml
```

Where `config.yml` is the name of your config file. Once you see confirmation that everything ran successfully, head over to Geckoboard and [start using your new Dataset to build widgets](https://support.geckoboard.com/hc/en-us/articles/223190488-Guide-to-using-datasets)!

## Building your config file

Here's what an example config file looks like:

```yaml
geckoboard_api_key: your_api_key
database:
 driver: mysql
 host: xxxx
 port: xxxx
 username: xxxx
 password: xxxx
 name: xxxx
 tls_config:
  ca_file: xxxx
  key_file: xxxx
  cert_file: xxxx
  ssl_mode: xxxx
refresh_time_sec: 60
datasets:
 - name: dataset.name
   update_type: replace
   sql: >
    SELECT 1, 0.34, source
    FROM table
   fields:
    - type: number
      name: Signups
    - type: percentage
      name: Conversion rate
    - type: string
      name: Source
```

#### Environment variables

If you wish, you can provide any of `geckoboard_api_key`, `host`, `port`, `username`, `password` and (database) `name` as environment variables with the syntax `"{{ YOUR_CUSTOM_ENV }}"`. Make sure to keep the quotes in there! For example:

```yaml
geckoboard_api_key: "{{ GB_API_KEY }}"
```

### geckoboard_api_key

Hopefully this is obvious, but this is where your Geckoboard API key goes. You can find yours [here](https://app.geckoboard.com/account/details).

### database

Enter the type of database you're connecting to in the `driver` field. SQL-Dataset supports:

- `mssql`
- `mysql`
- `postgres`
- `sqlite`

If you'd like to see support for another type of database, please raise a [support ticket](https://support.geckoboard.com/hc/en-us/requests/new?ticket_form_id=39437) or, if you're technically inclined, make the change and submit a pull request!

Only three parameters are required:

- `driver`
- `username`
- `name`

The other attributes, such as `host` and `port`, will default to their driver-specific values unless overridden.

#### SSL

If your database requires a CA cert or a x509 key/cert pair, you can supply this in `tls_config` under the database key.

```yaml
tls_config:
 ca_file: /path/to/file.pem
 key_file: /path/to/file.key
 cert_file: /path/to/cert.crt
 ssl_mode: (optional)
```

The possible values for `ssl_mode` depend on the database you're using:

- MSSQL: `disable`, `false`, `true` - try disable option if you experience connection issues
- MySQL: `true`, `skip-verify`
- Postgres: `disable`, `require`, `verify-ca`, `verify-full`
- SQLite: N/A


#### A note on user permissions

We _strongly_ recommend that the user account you use with SQL-Dataset has the lowest level of permission necessary. For example, one which is only permitted to perform `SELECT` statements on the tables you're going to be using. Like any SQL program, SQL-Dataset will run any query you give it, which includes destructive operations such as overwriting existing data, removing records, and dropping tables. We accept no responsibility for any adverse changes to your database due to accidentally running such a query.

### refresh_time_sec

Once started, SQL-Dataset can run your queries periodically and push the results to Geckoboard. Use this field to specify the time, in seconds, between refreshes.

If you do not wish for SQL-Dataset to run on a schedule, omit this option from your config.

### datasets

Here's where the magic happens - specify the SQL queries you want to run, and the Datasets you want to push their results into.

 - `name`: The name of your Dataset
 - `sql`: Your SQL query
 - `fields`: The schema of the Dataset into which the results of your SQL query will be parsed
 - `update_type`: Either `replace`, which overwrites the contents of the Dataset with new data on each update, or `append`, which merges the latest update with your existing data.
  - `unique_by`: An optional array of one or more field names whose values will be unique across all your records. When using the `append` update method, the fields in `unique_by` will be used to determine whether new data should update any existing records.

#### fields

A Dataset can hold up to 10 fields. The fields you declare should map directly to the columns that result from your `SELECT` query, in the **same order**.

For example:

```yaml
sql: SELECT date, orders, refunds FROM sales
fields:
 - name: Date
   type: date
 - name: Orders
   type: number
 - name: Refunds
   type: number
```

SQL-Dataset supports all of the field types supported by the [Datasets API](https://developer.geckoboard.com):

- date
- datetime
- duration
- number
- percentage
- string
- money

The `money` field type requires a `currency_code` to be provided:

```yaml
fields
 - name: MRR
   type: money
   currency_code: USD
```

The `duration` field type requires a `time_unit` to be provided:
With a value one of: milliseconds, seconds, minutes, hours

```yaml
fields
 - name: Time until support ticket resolved
   type: duration
   time_unit: minutes
```

Numeric field types can support null values. For a field to support this, pass the `optional` key:

```yaml
fields:
 - name: A field which might be NULL
   type: number
   optional: true
```

The Datasets API requires both a `name` and a `key` for each field, but SQL-Dataset will infer a `key` for you. Sometimes, however, the inferred `key` might not be permitted by the API. If you encounter such a case, you can supply a specific `key` value for that field.

```yaml
fields:
 - name: Your awesome field
   key: some_unique_key
   type: number
```
