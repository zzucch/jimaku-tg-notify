import jimaku_tg_notify/server
import dot_env as dot

pub fn main() {
  dot.new()
  |> dot.set_path(".env")
  |> dot.set_debug(False)
  |> dot.load

  server.run()
}
