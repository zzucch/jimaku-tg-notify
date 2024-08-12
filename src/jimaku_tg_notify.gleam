import dot_env as dot
import jimaku_tg_notify/http
import jimaku_tg_notify/server

pub fn main() {
  dot.new()
  |> dot.set_path(".env")
  |> dot.set_debug(False)
  |> dot.load

  let _ = http.send_message(123, "any")

  server.run()
}
