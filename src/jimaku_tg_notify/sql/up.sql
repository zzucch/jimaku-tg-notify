create table if not exists users (
  id integer primary key,
  telegram_chat_id integer
);

create table if not exists subscriptions (
  title_id primary key,
  user_id integer,
  latest_subtitile_time text,
  foreign key(user_id) references users(id)
);
