--SELECT b.id, s.name, b.title, b.num FROM books b, collections, series s
--WHERE fk_user = 1 AND fk_book = b.id AND fk_serie = s.id
--GROUP BY s.id, b.num, b.title, b.id ORDER BY s.name

--UPDATE users SET pwhash = crypt('bar', gen_salt('bf'))
--WHERE COALESCE(pwhash = crypt('baz', pwhash), FALSE) AND name = 'foo'

--WITH s AS (
--  SELECT id, activate, false FROM users WHERE email = 'susseppyrisa-7817@yopmail.com'
--), i AS (
--  INSERT INTO users ("email", "pwhash", "activate")
--  SELECT 'susseppyrisa-7817@yopmail.com', crypt('pkpas2ss3', gen_salt('bf')), gen_random_uuid()
--  WHERE NOT EXISTS (SELECT 1 FROM s)
--  RETURNING id, activate, true
--)
--SELECT id, activate, bool FROM i UNION ALL SELECT id, activate, bool FROM s

-- arnold = 9781849839730
-- fullmetal = 2351420187

--SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name, tags
--FROM books b
--INNER JOIN series s ON (b.fk_serie = s.id)
--LEFT OUTER JOIN collections ON (b.id = fk_book AND fk_user = 1)
--LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
--LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
----WHERE b.isbn = '2351420187'
--WHERE b.isbn = '9781849839730'

SELECT COUNT(*) FROM collections WHERE fk_user = 1 AND 'wishlist'=ANY(tags)
