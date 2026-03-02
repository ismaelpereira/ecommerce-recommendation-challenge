# The Challenge

This service was very interesting to do. I made it thinking on what Mahendra told me on the technical Interview. We have one endpoint to get all the events and we must to send. I was imagining the third party e-commerce was calling the POST /events and them we consume the other things.

If you follow the git commits, you will see a little bit of my strategy. I started with the simple and them I start to grow the project using the documentation I had access.

## Bigquery Decisions

I first created a new GCP project and setup a _gonear-test_ project. Them I created a Dataset called `ecommerce_events` and a table called `events`.

The table schema was the following

| Field      | Type      | Mode     | Default Value     |
| ---------- | --------- | -------- | ----------------- |
| event_id   | STRING    | REQUIRED | GENERATE_UUID     |
| user_id    | STRING    | REQUIRED | -                 |
| product_id | STRING    | REQUIRED | -                 |
| store_id   | STRING    | REQUIRED | -                 |
| event_type | STRING    | REQUIRED | -                 |
| timestamp  | TIMESTAMP | REQUIRED | CURRENT_TIMESTAMP |

- event_id -> I put a UUID because usualy I used sequencial numbers to be fast to read, but on the documentation nothing like this was there. I am generating this on the application level too so we can easily return on the request with the 201 code.

- timestamp -> I put as a default value the current date, but on the request is necessary so I also treat this on the code

### Partitioning

To have a strong Big Query Usage, I made a partition by date per day, I choose per day becausa it will be faster. I read that Big Query is optimized to have daily partitions and with basic math we can find this:

- 24 hours → scans 1–2 partitions
- 3 hours → scans 1 partition
- 36 hours → scans 2 partitions

And if we read 2 years of data. By daily Partitions we will scan **730** partitions (365x2). If we chose to partion Hourly (24*365*2) we will scan **17,520** partitions

### Clustering

I made a cluster by `store_id` and `product_id`. I understood it's a similar to the PostgreSQL index.

basic instead if we have

```
store_1 prod_7
store_8 prod_2
store_3 prod_9
store_1 prod_2
store_7 prod_1
```

with clustering we will have

```
store_1 prod_2
store_1 prod_7
store_1 prod_9
store_2 prod_1
store_2 prod_8
store_3 prod_4
```

So this are going to reduce the cost and the query reduction time

## Bigtable Decisions

The key for the row on Bigtable is defined as `user#{user_id}#revts#{reverse_timestamp}#evt%s`.

Where:

- `user_id` → ensures events are grouped by user
- `revts` → Reverse timestamp `math.MaxInt64 - event.Timestamp.UnixNano()`
- `event_id` → guarantees uniqueness and enables idempotency

If we are selecing one specific row we are going to get O(1) complexity and if we are searching for a event on specif time O(n) with n being the amount of user events since it will group all user events together

**Why Include `event_id`?**

Use only `user#{user_id}#revts#{reverse_timestamp}` will create a _data-loss risk_. 

If two events from the same user share the same timestamp, the second event that will arrive will overwrite the first causing a data-loss

With this format we guarante Row Key uniqueness, no same-time write overwrite, safe retry policy using `event_id` as idempotency key and a correlation with Big Query.

**Why Reverse Timestamp?**

Reverse timestamp will guarantee the last events appears first. This change the structure to a LIFO, instead of a FIFO too alowing more efficience on retrieving the last event without scanning the entire table.

### Tradeoffs 

This schema does not support efficient queries by `product_id`, `store_id` or `event_type`. But for this specific case, we are using the Big Query to agreggate ths information

## Error Handling strategy

I choose to treat most of the errors and for the query parameters explicity send the errors, not use default value for the amount of hours in case it's empty.

This causes 2 problems:

- On the error we explicity say the query parameter we are using (can be security issue)
- We will need to have a Good documentation, the API will be truncated for the user in case of setup.

**What happens if the BigQuery write succeeds but BigTable write fails?**

With what we have now, we are going to have inconsistent information, in case of failure. The current implementation performs a synchrounos dual-write. First it writes on the Big Query and them on Big Table. If Big Query succeesds and big table fails, the system will return an error and the information will be already stored on the table. If the user retry the request, we will have duplicated event on the Big Query table.

In production I will add a event ingestion via Pub/Sub, write only on the Big Query and send to Big Table trought a message queue. In this case Big Query is our source of truth.

To avoid problems we can add dead letter queue, with strong retry policy and send the event_id trought the Pub/Sub and use as a idempotency key, so we will not process the same event twice.

## What you’d add with more time

If I have more time I will Implement this as a separated API, Dockerize everything for local run, with a strong `docker-compose.yml` and deeper tests (This is something I already mention I need to improve with Golang).

To be production ready we can also add a cloud run deploy with CI/CD, but this will have costs to use Big Table on the GCP Project.

I think I will change Gin Gonic to the Framework you guys use too.

I wrote this before write tests, now checking how can I write the services unit tests I will create the interfaces like I did with the service, so we can easily mock the functions and write more tests.

Since tests are something I need to improve in Golang, I will write create a integration test to be more complete and safer to go to production.

I will add a package errors with a huge errors events and codes, so on the handlers we can treat better the errors.

## Environment Variables

```bash
# Project Config
PROJECT_ID=

# Bigtable Config
BIGTABLE_INSTANCE=
BIGTABLE_TABLE=
BIGTABLE_FAMILY=

#Bigquery Config
BIGQUERY_DATASET=
BIGQUERY_TABLE=
```

Right now I made a structure like this, because we only have one single table for Big Query and Big Table, on a production environment I would put something like BIGTABLE\_<RESOURCE>\_TABLE, so in this case, BIGTABLE_EVENTS_TABLE to be more explicit, but since we are going to work with Microsservices, it will be good to split trought different domains each service to not create a lot of responsability and avoid grow to much this environemnts variable to make the deploy easier and maintainable.

## Small Conclusion

I really liked the test, found very interesting the use case. I could see this analytics part of an ecommerce is way bigger than I tougth.

I learned the basic on how to implement Big Query and Big Table and why each one is important.
