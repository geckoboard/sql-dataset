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

As a starting point download the following [example config](docs/example.yml) - from here remove attributes you won't need and update the others.

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
