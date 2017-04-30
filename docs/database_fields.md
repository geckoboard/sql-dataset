## Database Examples

SQL dataset supports the following database drivers 

 - sqlite
 - mysql
 - postgres

* If you require an additional driver, please raise a [support ticket](https://support.geckoboard.com/hc/en-us/requests/new?ticket_form_id=39437)

Simplest database config requires at minimal just three keys, other attributes such as host and port default to specific driver defaults.

- driver 
- username
- name

#### Database user account

Please ensure to setup a user account has the least permissions (eg. only SELECT permissions on the required tables.) As its possible to execute any query just like any SQL program, and we can't be held responsible for any adverse changes to your database due to an accidental query

#### SSL

If your database requests a ca cert or a x509 key/cert pair, this is possible to supply as tls_config under the database key (examples for mysql/postgres below). `ssl_mode` is optional and has driver specific options.

```yaml
tls_config:
 ca_file: /file/path.pem
 key_file: /file/path.key 
 cert_file: /file/cert.crt
 ssl_mode: (optional, and datbase specific optionals)
```

## Full config examples for each database driver below 

### MySQL


```yaml
database:
 driver: mysql
 host: (optional)
 protocol: (tcp, unix - optional)
 username: (required)
 password: (optional)
 name: (database name - required)
 tls_config: (remove if not required - optional)
  ca_file: (file path, optional)
  key_file: (file path, optional) 
  cert_file: (file path, optional)
  ssl_mode: (true, skip-verify - optional)
 params: (optional)
  - customKey: customValue
```

### Postgres

```yaml
database:
 driver: postgres
 host: (optional)
 protocol: (tcp, unix - optional)
 username: (required)
 password: (optional)
 name: (database name - required)
 tls_config: (remove if not required - optional)
  ca_file: (file path, optional)
  key_file: (file path, optional) 
  cert_file: (file path, optional)
  ssl_mode: (disable, require, verify-ca, verify-full - optional)
 params: (optional)
  - customKey: customValue
```

### SQLite

```yaml
database:
 driver: sqlite
 password: (optional if required)
 name: (full file path to sqlite db)
 params:
 - Version: 4 (example param)
```