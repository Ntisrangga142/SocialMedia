CREATE TABLE public.followers (
    account_id   INT NOT NULL REFERENCES public.accounts(id),
    follower_id  INT NOT NULL REFERENCES public.accounts(id),
    read         BOOLEAN,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP NULL,
    CONSTRAINT followers_pk PRIMARY KEY (account_id, follower_id)
);
