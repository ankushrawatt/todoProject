DROP TABLE todo;
CREATE TABLE todo(
                     userid TEXT NOT NULL,
                     date TEXT NOT NULL,
                     task TEXT NOT NULL,
                     des TEXT NOT NULL,
                     createdAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                     CONSTRAINT fk_userid
                         FOREIGN KEY(userid)
                             REFERENCES users(userid)
);