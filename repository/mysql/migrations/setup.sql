DROP table Users;

CREATE TABLE Users
(
    id              int primary key auto_increment,
    name            varchar(255) not null,
    phone_number    varchar(255) not null unique,
    hashed_password varchar(500) not null,
    created_at      datetime     not null default current_timestamp
);

#
INSERT INTO game_app_db.Users(name, phone_number)
# VALUES ('user1', '09191234567'),
#        ('user2', '09121234567'),
#        ('user3', '09197654321');
#
#
SELECT *
           #
FROM game_app_db.Users #
where phone_number = '09197654321';
#
#
SELECT *
           #
FROM game_app_db.Users;