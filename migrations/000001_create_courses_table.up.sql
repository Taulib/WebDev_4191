-- Filename: migrations/000001_createcourses_table.up.sql

CREATE TABLE IF NOT EXISTS courses (
    coursesID varchar NOT NULL,
    coursesName text NOT NULL,
    creditHours int NOT NULL,
    /* createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW() */
);