CREATE TABLE public.comments (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id  INT NOT NULL REFERENCES public.accounts(id),
    post_id     INT NOT NULL REFERENCES public.posts(id),
    comment     TEXT,
    read        BOOLEAN,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL,
    deleted_at  TIMESTAMP NULL
);