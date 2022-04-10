CREATE TABLE `users` (
     `username` varchar(128) NOT NULL PRIMARY KEY,
     `hashed_password` varchar(256) NOT NULL,
     `full_name` varchar(256) NOT NULL,
     `email` varchar(128) NOT NULL,
     `password_changed_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);