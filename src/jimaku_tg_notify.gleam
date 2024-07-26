import gleam/hackney
import gleam/http/request
import gleam/int
import gleam/io
import gleam/list
import gleam/option.{None, Some}
import gleam/result.{try}
import gleam/string
import gleeunit/should
import nibble
import nibble/lexer

pub fn main() {
  let assert Ok(request) = request.to("https://jimaku.cc/entry/1563")

  use response <- try(
    request
    |> hackney.send,
  )

  response.status
  |> should.equal(200)

  io.debug(
    "response body length: "
    <> response.body
    |> string.length()
    |> int.to_string(),
  )

  let dates =
    response.body
    |> string.split(on: "\n")
    |> get_dates()

  io.debug(
    "dates amount: "
    <> dates
    |> list.length()
    |> int.to_string(),
  )

  Ok(dates)
}

pub type Date {
  Date(year: Int, month: Int, day: Int, hour: Int, minute: Int, second: Int)
}

pub type Token {
  Start
  End
  Plus
  Separator
  Colon
  Minus
  Num(Int)
}

fn lexer() {
  lexer.simple([
    lexer.token("<span class=\"table-data file-modified\" title=\"", Start),
    lexer.token("</span>", End),
    lexer.token("+", Plus),
    lexer.token("\">", Separator),
    lexer.token(":", Colon),
    lexer.token("-", Minus),
    lexer.int(Num),
    lexer.whitespace(Nil)
      |> lexer.ignore,
  ])
}

fn parser() {
  use _ <- nibble.do(nibble.token(Start))
  use year <- nibble.do(int_parser())
  use _ <- nibble.do(nibble.token(Minus))
  use month <- nibble.do(int_parser())
  use _ <- nibble.do(nibble.token(Minus))
  use day <- nibble.do(int_parser())
  use hour <- nibble.do(int_parser())
  use _ <- nibble.do(nibble.token(Colon))
  use minute <- nibble.do(int_parser())
  use _ <- nibble.do(nibble.token(Colon))
  use second <- nibble.do(int_parser())
  use _ <- nibble.do(offset_sign_parser())
  use _offset_hours <- nibble.do(int_parser())
  use _ <- nibble.do(nibble.token(Colon))
  use _offset_minutes <- nibble.do(int_parser())

  nibble.return(Date(year, month, day, hour, minute, second))
}

fn int_parser() {
  use token <- nibble.take_map("expected number")
  case token {
    Num(n) -> Some(n)
    _ -> None
  }
}

fn offset_sign_parser() {
  use token <- nibble.take_map("expected + or -")
  case token {
    Plus -> Some(Plus)
    Minus -> Some(Minus)
    _ -> None
  }
}

fn get_dates(data: List(String)) {
  list.filter_map(data, fn(line: String) {
    use tokens <- result.try(result.nil_error(lexer.run(line, lexer())))
    use date <- result.try(result.nil_error(nibble.run(tokens, parser())))
    Ok(date)
  })
}
