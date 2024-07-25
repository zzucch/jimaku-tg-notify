import gleam/hackney
import gleam/http/request
import gleam/io
import gleam/result.{try}
import gleam/string
import gleeunit/should

pub fn main() {
  let assert Ok(request) = request.to("https://jimaku.cc/entry/1563")

  use response <- try(
    request
    |> hackney.send,
  )

  response.status
  |> should.equal(200)

  io.debug(response.body)
  io.debug(string.length(response.body))

  Ok(response)
}
