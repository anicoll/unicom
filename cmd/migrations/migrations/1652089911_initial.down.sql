BEGIN;

DROP TABLE IF EXISTS response_channel;
DROP TABLE IF EXISTS response_channel_status;
DROP INDEX IF EXISTS communication;
DROP TYPE IF EXISTS communication_status;

COMMIT;