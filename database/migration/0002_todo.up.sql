DROP TABLE users;
CREATE TABLE users(
                      id UUID,
                      email varchar(50) NOT NULL,
                      password TEXT NOT NULL ,
                      userid Varchar(25) NOT NULL PRIMARY KEY,
                      firstname varchar(20) NOT NULL,
                      lastname varchar(20)NOT NULL,
                      mobile varchar(19) NOT NULL,
                      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                      UNIQUE (id,email,mobile)
);