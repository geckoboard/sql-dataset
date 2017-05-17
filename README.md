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

SQL-Datasets works by reading all of the information it needs from a YAML file. We've prepared an [example one](docs/example.yml) for you so you can get started quickly. The fields are fairly self-explanatory, but you can learn more about them below.

#### geckoboard_api_key

Hopefully this is obvious, but this is where your Geckoboard API key goes. You can find yours [here](https://app.geckoboard.com/account/details).

#### database

Enter the type of database you're connecting to in the `driver` field. SQL-Dataset supports:

- `mysql`
- `postgres`
- `sqlite`

If you'd like to see support for another type of database, please raise a [support ticket](https://support.geckoboard.com/hc/en-us/requests/new?ticket_form_id=39437) or, if you're technically inclined, make the change and submit a pull request!

Only three parameters are required:

- `driver`
- `username`
- `name`

The other attributes, such as `host` and `port`, will default to their driver-specific values unless overridden.

**A note on user accounts** - we _strongly_ recommend that the user account you use with SQL-Dataset has the lowest level of permission necessary. For example, a user account which is only permitted to perform `SELECT` statements on the tables you're going to be using. Like any SQL program, SQL-Dataset will run any query you give it, which includes destructive operations such as overwriting existing data, removing records, and dropping tables. We can't be held responsible for any adverse changes to your database due to accidentally running such a query.

If your database requires a CA cert or a x509 key/cert pair, you can supply this in `tls_config` under the database key.

```yaml
tls_config:
 ca_file: /path/to/file.pem
 key_file: /path/to/file.key 
 cert_file: /path/to/cert.crt
 ssl_mode: (optional)
```

The possible values for `ssl_mode` depend on the database you're using:

- MySQL: `true`, `skip-verify`
- Postgres: `disable`, `require`, `verify-ca`, `verify-full`
- SQLite: N/A

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
