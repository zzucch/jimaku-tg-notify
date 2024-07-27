import dot_env as dot
import gleam/int
import gleam/io
import gleam/list
import gleam/result
import gleam/string
import jimaku_tg_notify/http
import jimaku_tg_notify/parsing
import logging

pub fn main() {
  logging.configure()
  logging.set_level(logging.Debug)

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

  dot.new()
  |> dot.set_path(".env")
  |> dot.set_debug(False)
  |> dot.load
}
