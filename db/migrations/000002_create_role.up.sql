CREATE SCHEMA IF NOT EXISTS "role";

CREATE TABLE IF NOT EXISTS role.role_permissions (
    role varchar primary key,
    permissionID int
);

CREATE TABLE IF NOT EXISTS role.permissions (
    id serial primary key,
    sub_features int [] not null
);

CREATE TABLE IF NOT EXISTS role.sub_features (
    id serial primary key,
    featureID int not null,
    sub_feature_name varchar not null unique
);

CREATE TABLE IF NOT EXISTS role.features (
    id serial primary key,
    feature_name varchar unique not null
);

ALTER TABLE role.sub_features ADD FOREIGN KEY (featureID) REFERENCES role.features (id);
ALTER TABLE role.role_permissions ADD FOREIGN KEY (permissionID) REFERENCES role.permissions (id);

INSERT INTO role.features (feature_name) VALUES ('reminder'), ('dashboard'), ('backoffice');
INSERT INTO role.sub_features (sub_feature_name, featureID) VALUES ('allReminder', 1), ('allDashboard', 2), ('allBackoffice', 3);
INSERT INTO role.permissions (sub_features) VALUES ('{1,2,3}'),('{1}');
INSERT INTO role.role_permissions (role, permissionID) VALUES ('admin', 1), ('user', 2);