create table if not exists threads (
    threadID int not null,
    nextID int default 0,
    primary key (threadID)
);

create table if not exists comments (
    threadID int references threads,
    commentID int not null,
    author varchar(20),
    parentComment int,
    body varchar(400),
    primary key (threadID, commentID)
);

