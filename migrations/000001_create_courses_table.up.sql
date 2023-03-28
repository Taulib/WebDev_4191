-- Filename: migrations/000001_createcourses_table.up.sql

CREATE TABLE IF NOT EXISTS coursesT (
    coursesID varchar NOT NULL,
    coursesName text NOT NULL,
    creditHours int NOT NULL
);