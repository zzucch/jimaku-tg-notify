import gleam/dynamic
import gleam/int
import gleam/result
import simplifile
import sqlight

const data_dir = "./_data"

const connection = "file:./_data/sqlite.db?mode=rw"

const sql_path = "./src/jimaku_tg_notify/sql/"

pub type ErrorCode {
  UserNotFound
  FailedToFindUser
  FailedToSubscribe
  MultipleUsersWithOneChatId
}

pub fn start() {
  let assert Ok(Nil) = case simplifile.is_directory(data_dir) {
    Ok(True) -> Ok(Nil)
    _ -> simplifile.create_directory_all(data_dir)
  }
  use conn <- sqlight.with_connection(connection <> "c")
  let assert Ok(sql) = simplifile.read(sql_path <> "up.sql")
  sqlight.exec(sql, conn)
}

pub fn add_user(telegram_chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let assert Ok(sql) = simplifile.read(sql_path <> "add_user.sql")
  let sql = sql <> "(" <> int.to_string(telegram_chat_id) <> ");"
  sqlight.exec(sql, conn)
}

pub fn subscribe(
  telegram_chat_id: Int,
  title_id: Int,
  latest_subtitle_time: String,
) {
  use conn <- sqlight.with_connection(connection)
  use user_id <- result.try(find_user_id(telegram_chat_id))
  let assert Ok(sql) = simplifile.read(sql_path <> "subscribe.sql")
  let sql =
    sql
    <> "("
    <> int.to_string(title_id)
    <> ", "
    <> int.to_string(user_id)
    <> ", "
    <> latest_subtitle_time
    <> ");"
  case sqlight.exec(sql, conn) {
    Ok(Nil) -> Ok(Nil)
    _ -> Error(UserNotFound)
  }
}

pub fn unsubscribe(
  telegram_chat_id: Int,
  title_id: Int,
) {
  use conn <- sqlight.with_connection(connection)
  use user_id <- result.try(find_user_id(telegram_chat_id))
  let assert Ok(sql) = simplifile.read(sql_path <> "unsubscribe.sql")
  let sql =
    sql
    <> "("
    <> int.to_string(title_id)
    <> ", "
    <> int.to_string(user_id)
    <> ");"
  case sqlight.exec(sql, conn) {
    Ok(Nil) -> Ok(Nil)
    _ -> Error(UserNotFound)
  }
}

pub fn get_all_subscriptions(telegram_chat_id: Int) {
  todo
}

fn find_user_id(telegram_chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let decoder = dynamic.int
  let assert Ok(sql) = simplifile.read(sql_path <> "find_user.sql")
  let users =
    sqlight.query(
      sql,
      on: conn,
      with: [sqlight.int(telegram_chat_id)],
      expecting: decoder,
    )
  case users {
    Ok([first]) -> Ok(first)
    Ok([]) -> Error(UserNotFound)
    Ok(_) -> Error(MultipleUsersWithOneChatId)
    _ -> Error(FailedToFindUser)
  }
}
