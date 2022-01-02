CREATE TABLE players(
   playerID INTEGER PRIMARY KEY,
   name TEXT,
   actions TEXT,
   items TEXT
);

CREATE TABLE npcs(
   npcID INTEGER PRIMARY KEY,
   name TEXT NOT NULL,
   description TEXT,
   actions TEXT,
   dialog INTEGER
);

CREATE TABLE locations(
   locationName TEXT PRIMARY KEY,
   description  TEXT,
   actions TEXT,
   items TEXT,
   linkedLocations TEXT
);

CREATE TABLE items(
   itemName TEXT PRIMARY KEY,
   description  TEXT,
   actions TEXT
);

CREATE TABLE dialogs(
   dialogID INTEGER PRIMARY KEY,
   dialog  TEXT,
   responses TEXT,
   actions TEXT
);

CREATE TABLE actions (
   actionID INTEGER PRIMARY KEY,
   name TEXT,
   args TEXT
);


-- INSERT INTO items (name, description, actions) VALUES(
--    'Metal Trash Can Lid', 
--    'It stinks and has some nasty sludge on it... you might be able to hit someone with it.',
--    '"[1]"'
-- );

-- INSERT INTO players (name, actions, items) VALUES(
--    'Aesop', 
--    '"[2, 3]"',
--    '"[3, 4, 5]"'
-- );

-- INSERT INTO dialogs (dialog, responses, actions) VALUES(
--    'Scram!', 
--    '"[2]"',
--    '"[1]"'
-- );

-- INSERT INTO actions (name, args) VALUES(
--    'hit',
--    '"[\"metal trashcan lid\"]"' 
-- );

-- INSERT INTO npcs (name, description, actions) VALUES(
--    'Oscar', 
--    'He is a grouch because he lives in a trash can and people call him a grouch instead of helping him',
--    '"[]"'
-- );



