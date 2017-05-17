# Geckoboard SQL-Dataset

Quickly and easily send data from MySQL, Postgres and SQLite databases to Geckoboard Datasets.

SQL-Dataset is a command line app that takes the hassle out of integrating your database with Geckoboard. Rather than having to work with client libraries and write a bunch of code to connect to and query your database, with SQL-Dataset all you need to do is fill out a simple config file.

SQL-Dataset is available for macOS, Linux, and Windows. 

## How to use SQL-Dataset

### 1. Download the app

* [macOS](docs/macosx_setup.md)
* [Linux](docs/linux_setup.md)
* [Windows](docs/windows_setup.md)

### 2. Create a config file

SQL-Datasets works by reading all of the information it needs from a YAML file. We've prepared an [example one](docs/example.yml) for you so you can get started quickly.

#### geckoboard_api_key

Hopefully this is obvious, but this is where your Geckoboard API key goes. You can find yours [here](https://app.geckoboard.com/account/details).

#### database

```yml
database:
 driver: mysql
 host: 
 port: 3306
 username: jon_n
 password: my_password
 name: database_name
 tls_config:
  ca_file: path/to/file
  key_file: path/to/file
  cert_file: path/to/file
  ssl_mode: xxxx
```

The `driver` can be either `mysql`, `postgres` or `sqlite`.

`host`, `port`, `username`, `password` and `name` are the credentials you use to connect to your database.

`tls_config` is where you specify... [JON HELP].

#### refresh_time_sec

Once started, SQL-Dataset will run your queries periodically and push the results to Geckoboard. Use this field to specify the time, in seconds, between refreshes.

#### datasets

```
datasets:
 - name: dataset.name
   update_type: replace
   sql: >
    SELECT 1, 0.34, string
    FROM table
   fields:
    - type: number
      name: Count
    - type: percentage
      name: Some Percent
    - type: string
      name: Some Label
```

Here's where you specify the SQL queries you want to run, and the Datasets you want to push their results into.

Below are some references for the different dataset fields and database configurations for sqlite, postgres and mysql.

#### Top Level config options

Apart from the database and dataset top level keys there are two more one liners

##### Refresh Interval

The refresh interval describes whether the program should act as a scheduled process, and repeat the query at a set interval. For example you might want to push data every 10 seconds to the dataset. This is possible with the option `refresh_time_sec: 10`

##### Geckoboard API Key

This is just you Geckoboard API key for which you can retrieve from your Account section once logged into your Geckoboard account.

##### [Database attributes](docs/database_fields.md)
##### [Dataset & Field attributes](docs/dataset_fields.md)


### Build the widget from the Dataset

Head over to Geckoboard, and

 - Click 'Add Widget', and select the Datasets integration.
 - In the pop-out panel that appears you should see your new dataset.
 - You can use this to build a widget showing your data.
 - This will auto update every x seconds based on the config key value `refresh_time_sec`
