CREATE TABLE public.likes (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id  INT NOT NULL REFERENCES public.accounts(id),
    post_id     INT NOT NULL REFERENCES public.posts(id),
    read        BOOLEAN,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT likes_unique UNIQUE (account_id, post_id)
);