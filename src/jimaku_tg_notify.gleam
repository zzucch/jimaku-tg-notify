import gleam/int
import gleam/io
import gleam/list
import gleam/result
import gleam/string
import jimaku_tg_notify/http
import jimaku_tg_notify/parsing
import jimaku_tg_notify/storage

pub fn main() {
  let assert Ok(dates) = {
    use response <- result.map(http.get_response("https://jimaku.cc/entry/1563"))
    string.split(response, "\n")
    |> parsing.get_dates()
  }

  io.debug(
    "dates amount: "
    <> dates
    |> list.length()
    |> int.to_string(),
  )

  let assert Ok(_) = storage.start()
  let assert Ok(_) = storage.add_user(1)
  let assert Ok(_) = storage.get_all_subscriptions(1)
  let assert Ok(_) = storage.subscribe(1, 1, "")
  let assert Ok(_) = storage.get_all_subscriptions(1)
  let assert Ok(_) = storage.unsubscribe(1, 1)
  let assert Ok(_) = storage.get_all_subscriptions(1)

  Ok(dates)
}
