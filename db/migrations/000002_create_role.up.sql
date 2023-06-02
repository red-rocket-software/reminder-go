CREATE SCHEMA IF NOT EXISTS "role";

DROP TABLE IF EXISTS role.role_permissions;
CREATE TABLE IF NOT EXISTS role.role_permissions (
    role varchar primary key,
    permissions int [] not null
);

DROP TABLE IF EXISTS role.permissions;
CREATE TABLE IF NOT EXISTS role.permissions (
    id serial primary key,
    sub_features int [] not null
);

DROP TABLE IF EXISTS role.sub_features;
CREATE TABLE IF NOT EXISTS role.sub_features (
    id serial primary key,
    featureID int not null,
    name varchar not null
);

DROP TABLE IF EXISTS role.features;
CREATE TABLE IF NOT EXISTS role.features (
    id serial primary key,
    feature_name varchar unique not null
);

ALTER TABLE role.sub_features ADD FOREIGN KEY (featureID) REFERENCES role.features (id);

INSERT INTO role.features (feature_name) VALUES ('reminder'), ('dashboard');
INSERT INTO role.sub_features (name, featureID) VALUES ('allReminder', 1), ('allDashboard', 2);
INSERT INTO role.permissions (sub_features) VALUES ('{1,2}'),('{1}');
INSERT INTO role.role_permissions (role, permissions) VALUES ('admin', '{1}'), ('user', '{2}');