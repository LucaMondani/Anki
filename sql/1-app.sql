\c app;
CREATE TABLE decks (
    Id SERIAL PRIMARY KEY ,
    Title TEXT,
    Description TEXT
);
CREATE TABLE cards (
    Id SERIAL PRIMARY KEY ,
    Front TEXT,
    Back TEXT,
    Deck_Id int REFERENCES decks(Id) ON DELETE CASCADE
);