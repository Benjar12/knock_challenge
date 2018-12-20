# Knock Code Challenge

## Inital Thoughts

### First Thought 

When I first read the challenge I instantly jumped to using spark SQL. In spark SQL this would have been just a few lines.

```
val df = spark.read.format("csv").option("header", "true").load("person.csv")

df.write.format('jdbc').options(
      url='jdbc:mysql://localhost/knock_test',
      driver='com.mysql.jdbc.Driver',
      dbtable='people',
      user='knock_user',
      password='knock_password').mode('append').save()
```

There are some issues here. 
1. The challenge required it to be written in go.
2. A requirement was to handle ambiguity in the data. Spark will blow up.

### Second Thought

I then thought I should create a struct to match person. For each csv line I would create a type Person and use an ORM to store it. This does not satisfy the requirment to handle ambiguity.

### Final Thought

Given the requirements listed bellow I will read line by line. As we go through each line we will try to figure out the type unless it is already a string. In the event that there are two different types in the same column I plan to make that column a string.

## Requirements

### Hard Requirements

1. Read a csv file
2. Understand the type of each column
3. Handle inconsistent data
4. Store it in a new table per file
5. Atleast one test case

### Nice to Have

1. Handle CSVs without headers
2. Handle timestamps
3. Detect malformed csv files


## Known Issues

1. If a string parses to int then parses to float we convert it to string
2. Reading everything twice from the file system is expensive
3. Does not handle bad table or field names properly
4. Max field length should be evaluated not assumed

## Running the code

All of the dependancies are checked in the vendor folder and managed by godep. It should work out of the
box.

### The easy way

If you have docker-compose installed it's pretty straight foward to start things up.

1. go get github.com/Benjar12/knock_challenge
2. cd $GOPATH/src/Benjar12/knock_challenge
3. docker-compose up -d; docker-compose logs -f
4. Run curl

```curl --request POST \
  --url http://localhost:3302/load_file \
  --header 'cache-control: no-cache' \
  --header 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  --form file=@PATH TO FILE.csv \
  --form tablename=sample
  ```

### Running things manually

If you want to run things locally there is only one aditional step.

1. go get github.com/Benjar12/knock_challenge
2. cd $GOPATH/src/Benjar12/knock_challenge
3. export DATA_SOURCE="mysql_user:mysql_password@tcp(127.0.0.1:3306)/mysql_db"
4. go run main.go
5. Run curl

```curl --request POST \
  --url http://localhost:3302/load_file \
  --header 'cache-control: no-cache' \
  --header 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  --form file=@PATH TO FILE.csv \
  --form tablename=sample
  ```

## Closing thoughts.

This was a fun challenge. It got me thinking about other problems teams are facing. There is a lot of
refactoring I'd like todo. That said getting something that works was my priority.