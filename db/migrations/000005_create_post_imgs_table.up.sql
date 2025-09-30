CREATE TABLE public.post_imgs (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    post_id     INT NOT NULL REFERENCES public.posts(id),
    img         VARCHAR(255),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);