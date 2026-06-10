use cornucopia::{CodegenSettings, Error};
use openssl::ssl::{SslConnector, SslMethod, SslVerifyMode};
use postgres_openssl::MakeTlsConnector;
use config::CONFIG;

fn main() -> Result<(), Error> {
    println!("cargo:rerun-if-changed=queries");
    println!("cargo:rerun-if-changed=../../schema/schema.sql");

    let db_url = CONFIG.database_url.clone();

    let mut builder = SslConnector::builder(SslMethod::tls())
        .expect("Failed to create SSL connector");
    
    // Neon uses trusted certs; if you hit cert errors, you can temporarily
    // set SslVerifyMode::NONE, but PEER is preferred
    builder.set_verify(SslVerifyMode::PEER);
    let connector = MakeTlsConnector::new(builder.build());

    let mut client = postgres::Client::connect(&db_url, connector)
        .expect("Failed to connect to database");

    cornucopia::generate_live(
        &mut client,
        "queries",
        Some("src/generated.rs"),
        CodegenSettings {
            is_async: true,
            derive_ser: true,
        },
    )?;

    Ok(())
}