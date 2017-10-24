--SELECT b.id, s.name, b.title, b.num FROM books b, collections, series s
--WHERE fk_user = 1 AND fk_book = b.id AND fk_serie = s.id
--GROUP BY s.id, b.num, b.title, b.id ORDER BY s.name

UPDATE users SET pwhash = crypt('bar', gen_salt('bf'))
WHERE COALESCE(pwhash = crypt('baz', pwhash), FALSE) AND name = 'foo'

