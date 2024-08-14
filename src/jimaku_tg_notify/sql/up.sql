create table if not exists users (
  chat_id primary key
);

create table if not exists subscriptions (
  title_id primary key,
  chat_id integer,
  latest_subtitle_time integer,
  foreign key(chat_id) references users(chat_id)
);
