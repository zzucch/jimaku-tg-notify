import gleam/int
import simplifile
import sqlight

const data_dir = "./_data"

const connection = "file:./_data/sqlite.db?mode=rw"

const sql_path = "./src/jimaku_tg_notify/sql/"

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

pub fn subscribe() {
  todo
}

pub fn unsubscribe() {
  todo
}

pub fn get_all_subscriptions() {
  todo
}
