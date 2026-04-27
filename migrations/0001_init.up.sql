BEGIN;

CREATE TABLE announcements (
    id BIGSERIAL PRIMARY KEY,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_by UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    title TEXT NOT NULL,
    description TEXT NOT NULL
);

COMMIT;