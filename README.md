Go JSONB Test
-------------

Sample project to test storing and retrieving binary value inside json object in MySQL using Golang.



### Steps:

- Run the app along with mysql in docker:

    ```shell script
    make run
    ```

or

- You can also the run the app in local (requires go 1.14) but start the db in docker:

    ```shell script
    make run-db
    make run-local
    ```

- It also has some tests, but it's similar to the main function (so not important):
    ```shell script
    make run-db
    make test
    ```

--------
**Note**:
- To change the database address and name, use the local.env or docker-compose.xml file.


#### Project Structure:

- When the app starts, it migrates the database using the sql in `./migrations` directory.
It inserts test data and reads it using different collation option.


