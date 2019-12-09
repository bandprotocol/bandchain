pub mod arc_module;
mod colors;

use arc_module::ArcModule;
use colors::Rgb;

#[no_mangle]
pub extern fn apply() {
  let mut module = ArcModule::get_instance();
  let rows = module.rows;
  let cols = module.cols;
  let ref mut animation = module.get_animation().as_mut_slice();
  for (index, frame) in animation.chunks_mut(rows * cols).enumerate() {
    for row in 0 .. rows {
      for col in 0 .. cols {
        frame[row * cols + col] = Rgb::new(row as u8 * 6, col as u8 * 6, ((index as u8 * 6) % 0xff));
      }
    }
  }
}
