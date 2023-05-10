-- Filename: migrations/000001_createcourses_table.up.sql

CREATE TABLE IF NOT EXISTS courses (
    CoursesID int NOT NULL,
    CourseName text NOT NULL,
    CreditHours int NOT NULL
);