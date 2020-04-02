CREATE TABLE IF NOT EXISTS t1
(
    id                        BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name                      VARCHAR(32),
    external_uuids            JSON,
#     external_uuids_binary     JSON COLLATE `binary`,
#     external_uuids_mb4_binary JSON COLLATE utf8mb4_bin, # doesn't allow for json,
    PRIMARY KEY (id)
#     ,INDEX external_uuids_idx ((CAST(external_uuids AS BINARY(16) ARRAY)))
)
    DEFAULT CHARACTER SET = 'utf8mb4',
    DEFAULT COLLATE = 'utf8mb4_unicode_ci'
;
