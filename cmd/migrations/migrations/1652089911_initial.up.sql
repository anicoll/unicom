BEGIN;

CREATE TYPE communication_status AS ENUM('PENDING', 'SUCCESS', 'FAILED');
CREATE TYPE notification_type AS ENUM('EMAIL', 'SMS', 'PUSH');

CREATE TABLE IF NOT EXISTS communications (
  id TEXT NOT NULL,
  domain TEXT NOT NULL,
  "type" notification_type NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "status" communication_status DEFAULT 'PENDING',
  PRIMARY KEY (id)
);

CREATE TYPE response_channel_status AS ENUM('PENDING', 'SUCCESS', 'FAILED');
CREATE TYPE response_channel_type AS ENUM('SQS', 'WEBHOOK', 'EVENT_BRIDGE');

CREATE TABLE IF NOT EXISTS response_channels (
  id TEXT NOT NULL,
  communication_id TEXT NOT NULL REFERENCES communications (id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  external_id TEXT,
  "status" response_channel_status NOT NULL DEFAULT 'PENDING',
  "type" response_channel_type NOT NULL,
  "url" TEXT NOT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_response_channel_communication_id ON response_channels (communication_id);

COMMIT;