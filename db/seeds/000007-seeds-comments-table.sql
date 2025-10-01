INSERT INTO public.comments (account_id, post_id, comment, read, created_at)
VALUES
(2, 1, 'Nice first post, Alice!', true, NOW()),
(3, 1, 'Welcome to the platform', true, NOW()),
(1, 2, 'That coffee looks delicious', true, NOW()),
(4, 3, 'Congrats on the run!', true, NOW()),
(5, 4, 'Wow, amazing view', true, NOW());
