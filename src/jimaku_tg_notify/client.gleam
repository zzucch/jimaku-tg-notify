import gleam/hackney
import gleam/http
import gleam/http/request
import gleam/int
import gleam/json
import gleam/list
import gleam/result
import gleam/string
import gleeunit/should
import jimaku_tg_notify/parsing
import logging

pub fn get_dates() {
  logging.configure()
  logging.set_level(logging.Debug)

  let assert Ok(dates) = {
    use response <- result.map(get_response("https://jimaku.cc/entry/1563"))
    string.split(response, "\n")
    |> parsing.get_dates()
  }

  logging.log(
    logging.Debug,
    "dates amount: "
      <> dates
    |> list.length()
    |> int.to_string(),
  )
}

pub fn send_message(chat_id: Int, message: String) {
  let assert Ok(request) = request.to("http://localhost:3002/sendMessage")

  use response <- result.try(
    request
    |> request.set_method(http.Post)
    |> request.set_header("Content-Type", "application/json")
    |> request.set_body(message_to_json(chat_id, message))
    |> hackney.send(),
  )

  response.status
  |> should.equal(200)

  Ok(response.body)
}

fn get_response(url: String) {
  let assert Ok(request) = request.to(url)

  use response <- result.try(
    request
    |> hackney.send,
  )

  response.status
  |> should.equal(200)

  Ok(response.body)
}

fn message_to_json(chat_id: Int, message: String) {
  json.object([#("chat_id", json.int(chat_id)), #("text", json.string(message))])
  |> json.to_string
}
