INSERT INTO public.accounts (email, password, created_at)
VALUES
('alice@example.com', '$argon2id$v=19$m=65536,t=2,p=1$aAgJzdCO1OCvabVhOF7quA$14w5enYmDzC4MxBKrIUyHZqTzq7Z3RG9h71fVMoaNsY', NOW()),
('bob@example.com', '$argon2id$v=19$m=65536,t=2,p=1$aAgJzdCO1OCvabVhOF7quA$14w5enYmDzC4MxBKrIUyHZqTzq7Z3RG9h71fVMoaNsY', NOW()),
('charlie@example.com', '$argon2id$v=19$m=65536,t=2,p=1$aAgJzdCO1OCvabVhOF7quA$14w5enYmDzC4MxBKrIUyHZqTzq7Z3RG9h71fVMoaNsY', NOW()),
('diana@example.com', '$argon2id$v=19$m=65536,t=2,p=1$aAgJzdCO1OCvabVhOF7quA$14w5enYmDzC4MxBKrIUyHZqTzq7Z3RG9h71fVMoaNsY', NOW()),
('eric@example.com', '$argon2id$v=19$m=65536,t=2,p=1$aAgJzdCO1OCvabVhOF7quA$14w5enYmDzC4MxBKrIUyHZqTzq7Z3RG9h71fVMoaNsY', NOW());
