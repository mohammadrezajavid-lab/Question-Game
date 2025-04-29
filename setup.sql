DROP table Users;

CREATE TABLE Users
(
    Id          int primary key auto_increment,
    Name        varchar(255) not null,
    PhoneNumber varchar(255) not null unique,
    CreatedAt   datetime     not null default current_timestamp
);

INSERT INTO game_app_db.Users(Name, PhoneNumber)
VALUES ('user1', '09191234567'),
       ('user2', '09121234567'),
       ('user3', '09197654321');