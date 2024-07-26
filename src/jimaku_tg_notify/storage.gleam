import simplifile
import sqlight

pub fn initialize_storage() {
  use conn <- sqlight.with_connection(":memory:")

  let assert Ok(sql) = simplifile.read("./src/jimaku_tg_notify/sql/up.sql")
  let assert Ok(Nil) = sqlight.exec(sql, conn)
}

pub fn add_user() {
  todo
}

pub fn subscribe() {
  todo
}

pub fn unsubscribe() {
  todo
}

pub fn check_all_subscriptions() {
  todo
}
