## Dataset Attributes

 - **Name:** The name of your dataset generally in the format of `name.part.c`
 - **SQL:** Your sql select query
 - **Unique By:** An optional array of one or more field names whose values will be unique across all your records.

**Update Type**
  - **Replace:** Sends the first 500 records and will error that it only sent the first 500 rows if your sql rows extend over that
  - **Append:** Sends all the records returned from the sql resultset. This supports upto 5000 rows and will remove older records.

```yaml
datasets:
 - name: Your database name
   update_type: (replace, append - optional)
   sql: SELECT column, count(*) FROM table GROUP BY column
   unique_by: ["field name"] (optional)

```

## Field Attributes

As per the [datasets api reference](https://developer.geckoboard.com/api-reference/curl/) sql-dataset supports the following field types.

A dataset supports upto **10 fields** and fields must be declared in the **same order** as your sql select query.

- date
- datetime
- number
- percentage
- string
- money

##### Example field entry for a dataset

```yaml
fields:
 - name: Your Label
   type: number
```

##### Multiple fields

```yaml
fields:
 - name: Your Label one
   type: string
 - name: Your Label two
   type: percentage
 - name: Your Label three
   type: datetime
```

##### Money field type requires a `currency code` key

```yaml
fields
 - name: MRR
   type: money
   currency_code: USD
```

##### Number field type can support null values to support this, pass `optional` key

```yaml
fields:
 - name: Your Label one
   type: number
   optional: true
```

### Field unique key

Some field names might be the same or will not be permitted as a key to Geckoboard, in these cases you can supply a specific key value for any field with the `key` attribute.

When the key exists this will be used instead of generating it from the name. This key must contain only letters, underscore and numbers (but not the first or last character)

```yaml
fields:
 - name: Your Label one
   key: some_unique_key
   type: number
```
