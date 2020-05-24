create table executions(
    id integer auto_increment primary key,
    created_at timestamp not null,
    seconds integer not null ,
    zone varchar(5) not null
);