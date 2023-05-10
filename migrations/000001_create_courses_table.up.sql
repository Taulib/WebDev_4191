-- Filename: migrations/000001_createcourses_table.up.sql

CREATE TABLE IF NOT EXISTS courses (
    CourseID bigserial PRIMARY KEY,
    CourseName text NOT NULL,
    CreditHours text NOT NULL
);