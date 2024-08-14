import gleam/dynamic
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
  FailedToGetSubscriptions
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

pub fn add_user(chat_id: Int) {
  use conn <- sqlight.with_connection(connection)
  let assert Ok(sql) = simplifile.read(sql_path <> "add_user.sql")
  sqlight.query(sql, conn, [sqlight.int(chat_id)], dynamic.dynamic)
}

pub fn subscribe(chat_id: Int, title_id: Int, latest_subtitle_time: Int) {
  use conn <- sqlight.with_connection(connection)
  let assert Ok(sql) = simplifile.read(sql_path <> "subscribe.sql")
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
  let assert Ok(sql) = simplifile.read(sql_path <> "unsubscribe.sql")
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
  let assert Ok(sql) = simplifile.read(sql_path <> "get_all_subscriptions.sql")
  let res = sqlight.query(sql, conn, [sqlight.int(chat_id)], decoder)
  case res {
    Ok(_) -> Ok(Nil)
    _ -> Error(FailedToGetSubscriptions)
  }
}
