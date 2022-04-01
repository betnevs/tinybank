CREATE TABLE `accounts` (
    `id` bigint PRIMARY KEY,
    `owner` varchar(128) NOT NULL,
    `balance` bigint NOT NULL,
    `currency` varchar(32) NOT NULL,
    `created_at` datetime NOT NULL
);

CREATE TABLE `entries` (
    `id` bigint PRIMARY KEY,
    `account_id` bigint NOT NULL,
    `amount` bigint NOT NULL,
    `created_at` datetime NOT NULL
);

CREATE TABLE `transfers` (
    `id` bigint PRIMARY KEY,
    `from_account_id` bigint NOT NULL,
    `to_account_id` bigint NOT NULL,
    `amount` bigint NOT NULL,
    `created_at` datetime NOT NULL
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);

CREATE INDEX `entries_index_1` ON `entries` (`account_id`);

CREATE INDEX `transfers_index_2` ON `transfers` (`from_account_id`, `to_account_id`);

CREATE INDEX `transfers_index_3` ON `transfers` (`to_account_id`);
