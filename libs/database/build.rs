
use cornucopia::{CodegenSettings, Error};

fn main() -> Result<(), Error> {
    println!("cargo:rerun-if-changed=queries");
    println!("cargo:rerun-if-changed=../../schema/schema.sql");

    cornucopia::generate_managed(
        "../../schema",
        vec![String::from("queries/**")],
        Some("src/generated.rs"),
        false,
        CodegenSettings {
           is_async: true,
           derive_ser: true 
        },
    )?;

    Ok(())
}