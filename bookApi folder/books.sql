CREATE DATABASE mydatabase;
use mydatabase;

CREATE TABLE mydatabase.books(
    ID INT NOT NULL,
    Title VARCHAR(45) NOT NULL,
    Author VARCHAR(45) NOT NULL,
    Price INT NOT NULL,
    PRIMARY KEY (ID)
);

INSERT INTO books(ID, Title, Author, Price)
VALUES(1, 'TEST BOOK 1', 'AUTHOR 1', 200);