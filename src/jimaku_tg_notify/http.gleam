import gleam/hackney
import gleam/http
import gleam/http/request
import gleam/json
import gleam/result
import gleeunit/should

pub fn get_response(url: String) {
  let assert Ok(request) = request.to(url)

  use response <- result.try(
    request
    |> hackney.send,
  )

  response.status
  |> should.equal(200)

  Ok(response.body)
}

pub fn send_message(chat_id: Int, message: String) {
  let assert Ok(request) = request.to("http://localhost:8080/sendMessage")

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

fn message_to_json(chat_id: Int, message: String) {
  json.object([#("chat_id", json.int(chat_id)), #("text", json.string(message))])
  |> json.to_string
}
