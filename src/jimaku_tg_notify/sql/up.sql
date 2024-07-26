create table if not exists user_data (
  id integer primary key,
  telegram_chat_id integer
);

create table if not exists subscription_data (
  title_id primary key,
  latest_subtitile_time text
);
