CREATE SCHEMA IF NOT EXISTS "role";

CREATE TABLE IF NOT EXISTS role.role_permissions (
    role varchar primary key,
    permissions int [] not null
);

CREATE TABLE IF NOT EXISTS role.permissions (
    id serial primary key,
    features int [] not null ,
    sub_features int [] not null
);

CREATE TABLE IF NOT EXISTS role.features (
    id serial primary key,
    feature_name varchar unique not null
);

CREATE TABLE IF NOT EXISTS role.sub_features (
    id serial primary key,
    featureID int not null,
    name varchar not null
);

ALTER TABLE role.sub_features ADD FOREIGN KEY (featureID) REFERENCES role.features (id);