INSERT INTO role.sub_features (name, featureID) VALUES ('all', 1), ('all', 2);
INSERT INTO role.features (feature_name) VALUES ('reminder'), ('dashboard');
INSERT INTO role.permissions (features, sub_features) VALUES ('{1,2}', '{1,2}');
INSERT INTO role.role_permissions (role, permissions) VALUES ('admin', '{1}');