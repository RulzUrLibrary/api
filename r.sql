--SELECT b.id, s.name, b.title, b.num FROM books b, collections, series s
--WHERE fk_user = 1 AND fk_book = b.id AND fk_serie = s.id
--GROUP BY s.id, b.num, b.title, b.id ORDER BY s.name

--UPDATE users SET pwhash = crypt('bar', gen_salt('bf'))
--WHERE COALESCE(pwhash = crypt('baz', pwhash), FALSE) AND name = 'foo'

WITH s AS (
  SELECT id, activate, false FROM users WHERE email = 'susseppyrisa-7817@yopmail.com'
), i AS (
  INSERT INTO users ("email", "pwhash", "activate")
  SELECT 'susseppyrisa-7817@yopmail.com', crypt('pkpas2ss3', gen_salt('bf')), gen_random_uuid()
  WHERE NOT EXISTS (SELECT 1 FROM s)
  RETURNING id, activate, true
)
SELECT id, activate, bool FROM i UNION ALL SELECT id, activate, bool FROM s
