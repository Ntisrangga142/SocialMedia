CREATE TABLE public.posts (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id  INT NOT NULL REFERENCES public.accounts(id),
    caption     TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);