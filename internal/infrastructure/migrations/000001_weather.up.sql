create table weather(
    id int auto_increment primary key ,
    weather_value int not null ,
    created_at timestamp not null ,
    type varchar(20) not null
);