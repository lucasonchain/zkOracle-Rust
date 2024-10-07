use clap::Parser;
use std::path::PathBuf;

mod simulator;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[arg(short = 'r', long, default_value_t = 10)]
    runs: u32,

    #[arg(short = 'd', long, default_value = "./data.csv")]
    dst: PathBuf,

    #[arg(short = 'e', long, default_value = "ws://127.0.0.1:8545/")]
    eth_url: String,

    #[arg(short = 'c', long, default_value = "0x40918Ba7f132E0aCba2CE4de4c4baF9BD2D7D849")]
    contract: String,

    #[arg(short = 'k', long, default_value = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")]
    private_key: String,

    #[arg(short = 'm', long, default_value_t = 0)]
    mode: u8,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();

    let aggregation_circuit_analyzer = simulator::Simulator::new(
        args.runs,
        args.dst,
        simulator::SimulationMode::from(args.mode),
        args.eth_url,
        args.contract,
        args.private_key,
    )?;

    aggregation_circuit_analyzer.simulate()?;

    Ok(())
}