import gleam/list
import gleam/option.{None, Some}
import gleam/result
import nibble
import nibble/lexer

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

pub fn get_dates(data: List(String)) {
  list.filter_map(data, fn(line: String) {
    use tokens <- result.try(result.nil_error(lexer.run(line, lexer())))
    use date <- result.try(result.nil_error(nibble.run(tokens, parser())))
    Ok(date)
  })
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
