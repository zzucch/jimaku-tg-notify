import gleam/bytes_builder
import gleam/erlang/process
import gleam/http/request.{type Request}
import gleam/http/response.{type Response}
import gleam/int
import jimaku_tg_notify/storage
import mist.{type Connection, type ResponseData}

pub fn run() {
  let assert Ok(Nil) = storage.start()
  let assert Ok(_) =
    fn(req: Request(Connection)) -> Response(ResponseData) {
      case request.path_segments(req) {
        ["sub", ..rest] -> handle_sub(rest)
        ["unsub", ..rest] -> handle_unsub(rest)
        _ -> not_found()
      }
    }
    |> mist.new
    |> mist.port(3001)
    |> mist.start_http

  process.sleep_forever()
}

fn bad_request() {
  response.new(400)
  |> response.set_body(mist.Bytes(bytes_builder.new()))
}

fn not_found() {
  response.new(404)
  |> response.set_body(mist.Bytes(bytes_builder.new()))
}

fn handle_sub(args: List(String)) {
  case args {
    [first, second] -> {
      let assert Ok(chat_id) = int.parse(first)
      let assert Ok(title_id) = int.parse(second)

      let _ = storage.add_user(chat_id)
      let assert Ok(_) = storage.subscribe(chat_id, title_id, 0)

      response.new(200)
      |> response.set_body(mist.Bytes(bytes_builder.new()))
    }
    _ -> bad_request()
  }
}

fn handle_unsub(args: List(String)) {
  case args {
    [first, second] -> {
      let assert Ok(chat_id) = int.parse(first)
      let assert Ok(title_id) = int.parse(second)

      let assert Ok(_) = storage.unsubscribe(chat_id, title_id)

      response.new(200)
      |> response.set_body(mist.Bytes(bytes_builder.new()))
    }
    _ -> bad_request()
  }
}
