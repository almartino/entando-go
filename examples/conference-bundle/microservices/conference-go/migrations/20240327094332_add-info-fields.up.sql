ALTER TABLE conferences
    ADD COLUMN create_ts timestamp,
    ADD COLUMN create_by char(36),
    ADD COLUMN update_ts timestamp,
    ADD COLUMN update_by char(36);
