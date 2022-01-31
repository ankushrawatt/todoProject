CREATE TABLE users(
    ID UUID,
    Email varchar(50) NOT NULL,
    Password TEXT NOT NULL ,
    UserId Varchar(25) NOT NULL ,
    FirstName varchar(20) NOT NULL,
    LastName varchar(20)NOT NULL,
    MobileNo varchar(19) NOT NULL,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);