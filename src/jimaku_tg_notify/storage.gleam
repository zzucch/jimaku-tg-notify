import gleam/dynamic
import simplifile
import sqlight

const data_dir = "./_data"

const connection = "file:./_data/sqlite.db?mode=rw"

pub type ErrorCode {
  UserNotFound
  FailedToFindUser
  FailedToSubscribe
  FailedToUnsubscribe
  FailedToGetSubscriptions
  MultipleUsersWithOneChatId
}

pub fn start() {
  let assert Ok(Nil) = case simplifile.is_directory(data_dir) {
    Ok(True) -> Ok(Nil)
    _ -> simplifile.create_directory_all(data_dir)
  }
  use conn <- sqlight.with_connection(connection <> "c")
  let sql =
    "create table if not exists users (
      chat_id primary key
    );

    create table if not exists subscriptions (
      title_id primary key,
      chat_id integer,
      latest_subtitle_time integer,
      foreign key(chat_id) references users(chat_id)
    );"
  sqlight.exec(sql, conn)
}

pub fn add_user(chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let sql = "insert into users (chat_id) values (?);"

  sqlight.query(sql, conn, [sqlight.int(chat_id)], dynamic.dynamic)
}

pub fn subscribe(chat_id: Int, title_id: Int, latest_subtitle_time: Int) {
  use conn <- sqlight.with_connection(connection)
  let sql =
    "insert into subscriptions (
       title_id,
       chat_id,
       latest_subtitle_time) 
    values (?, ?, ?);"
  case
    sqlight.query(
      sql,
      conn,
      [
        sqlight.int(title_id),
        sqlight.int(chat_id),
        sqlight.int(latest_subtitle_time),
      ],
      dynamic.dynamic,
    )
  {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToSubscribe)
  }
}

pub fn unsubscribe(chat_id: Int, title_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let sql = "delete from subscriptions where title_id = ? and chat_id = ?;"
  case
    sqlight.query(
      sql,
      conn,
      [sqlight.int(title_id), sqlight.int(chat_id)],
      dynamic.dynamic,
    )
  {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToUnsubscribe)
  }
}

pub fn get_all_subscriptions(chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let decoder = dynamic.element(0, dynamic.int)
  let sql =
    "select title_id, latest_subtitle_time from subscriptions 
     where chat_id = ?"
  let res = sqlight.query(sql, conn, [sqlight.int(chat_id)], decoder)
  case res {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToGetSubscriptions)
  }
}
