import gleam/int
import gleam/list
import gleam/result
import gleam/string
import jimaku_tg_notify/http
import jimaku_tg_notify/parsing
import logging

pub fn run() {
  logging.configure()
  logging.set_level(logging.Debug)

  let assert Ok(dates) = {
    use response <- result.map(http.get_response("https://jimaku.cc/entry/1563"))
    string.split(response, "\n")
    |> parsing.get_dates()
  }

  logging.log(
    logging.Critical,
    "dates amount: "
      <> dates
    |> list.length()
    |> int.to_string(),
  )
}
