use std::env;
use std::process;
use std::sync::Arc;
use tokio::signal;

mod config;
mod zk_oracle;

#[tokio::main]
async fn main() {
    let config_file = env::args()
        .nth(1)
        .unwrap_or_else(|| String::from("./configs/config.example.json"));

    let config = config::load_config(&config_file).unwrap_or_else(|err| {
        eprintln!("Failed to load config: {}", err);
        process::exit(1);
    });

    let node = Arc::new(zk_oracle::Node::new(config).unwrap_or_else(|err| {
        eprintln!("Failed to create node: {}", err);
        process::exit(1);
    }));

    let node_clone = Arc::clone(&node);
    tokio::spawn(async move {
        if let Err(e) = node_clone.run().await {
            eprintln!("Node error: {}", e);
            process::exit(1);
        }
    });

    signal::ctrl_c().await.expect("Failed to listen for ctrl+c");
    println!("Shutting down...");

    node.stop().await;
}