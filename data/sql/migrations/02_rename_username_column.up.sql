ALTER TABLE users RENAME COLUMN username TO email; 
ALTER TABLE users DROP COLUMN mnemonic, DROP COLUMN password;
Alter Table users ALTER COLUMN xpriv TYPE text;
