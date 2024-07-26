import gleam/dynamic
import gleam/int
import gleam/io
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
  FailedToUnsubscribe
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
  sqlight.query(sql, conn, [sqlight.int(telegram_chat_id)], dynamic.dynamic)
}

pub fn subscribe(
  telegram_chat_id: Int,
  title_id: Int,
  latest_subtitle_time: String,
) {
  use conn <- sqlight.with_connection(connection)
  use user_id <- result.try(find_user_id(telegram_chat_id))
  let assert Ok(sql) = simplifile.read(sql_path <> "subscribe.sql")
  case
    sqlight.query(
      sql,
      conn,
      [
        sqlight.int(title_id),
        sqlight.int(user_id),
        sqlight.text(latest_subtitle_time),
      ],
      dynamic.dynamic,
    )
  {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToSubscribe)
  }
}

pub fn unsubscribe(telegram_chat_id: Int, title_id: Int) {
  use conn <- sqlight.with_connection(connection)
  use user_id <- result.try(find_user_id(telegram_chat_id))
  let assert Ok(sql) = simplifile.read(sql_path <> "unsubscribe.sql")
  case
    sqlight.query(
      sql,
      conn,
      [sqlight.int(title_id), sqlight.int(user_id)],
      dynamic.dynamic,
    )
  {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToUnsubscribe)
  }
}

pub fn get_all_subscriptions(telegram_chat_id: Int) {
  todo
}

fn find_user_id(telegram_chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let decoder = dynamic.element(0, dynamic.int)
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
