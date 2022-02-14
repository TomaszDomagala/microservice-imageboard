create table if not exists boards
(
    id   varchar(4)  not null CHECK (length(id) > 0),
    name varchar(30) not null,
    primary key (id)
);

insert into boards (id, name)
values ('math', 'Math'),
       ('prog', 'Programming'),
       ('ck', 'Cooking'),
       ('s', 'Sports'),
       ('alc', 'Alcohol (Piwo)'),
       ('bg', 'Board games'),
       ('anim', 'hinskie bajki'),
       ('sci', 'Science'),
       ('gm', 'Gaming'),
       ('pol', 'Politics');

create table if not exists threads
(
    boardID    varchar(4)   not null references boards(id),
    threadID serial          not null,
    owner    varchar(255) not null
);
