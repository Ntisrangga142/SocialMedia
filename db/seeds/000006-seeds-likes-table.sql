INSERT INTO public.likes (account_id, post_id, read, created_at)
VALUES
(2, 1, true, NOW()),  -- Bob likes Alice's post
(3, 1, true, NOW()),  -- Charlie likes Alice's post
(1, 2, true, NOW()),  -- Alice likes Bob's post
(4, 2, true, NOW()),  -- Diana likes Bob's post
(5, 4, true, NOW());  -- Eric likes Diana's post

