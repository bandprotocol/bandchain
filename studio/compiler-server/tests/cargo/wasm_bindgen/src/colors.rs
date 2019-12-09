#[repr(C)]
#[derive(Copy, Clone)]
pub struct Rgb {
  pub r: u8,
  pub g: u8,
  pub b: u8,
}

impl Rgb {
  pub fn new(r: u8, g: u8, b: u8) -> Rgb {
    Rgb {r, g, b}
  }
}