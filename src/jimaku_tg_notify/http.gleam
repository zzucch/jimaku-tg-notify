import gleam/hackney
import gleam/http/request
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
