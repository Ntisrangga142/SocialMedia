CREATE TABLE public.profiles (
    id          INT PRIMARY KEY REFERENCES public.accounts(id),
    fullname    VARCHAR(255),
    phone       VARCHAR(255),
    img         VARCHAR(255),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NULL
);
