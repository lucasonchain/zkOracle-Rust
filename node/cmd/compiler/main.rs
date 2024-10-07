use std::env;
use std::fs::{self, File};
use std::io::Write;
use std::path::{Path, PathBuf};

mod zk_oracle;

// You'll need to find or create Rust equivalents for these
use zk_snark_library::{Circuit, R1CS, Groth16, BN254};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    let build_path = args.get(1).map(PathBuf::from).unwrap_or_else(|| PathBuf::from("./build/"));

    compile(&zk_oracle::AggregationCircuit::new(), build_path.join("aggregation"))?;
    compile(&zk_oracle::SlashingCircuit::new(), build_path.join("slashing"))?;

    Ok(())
}

fn compile(circuit: &dyn Circuit, dst: PathBuf) -> Result<(), Box<dyn std::error::Error>> {
    if !dst.exists() {
        fs::create_dir_all(&dst)?;
    }

    let r1cs = R1CS::compile(BN254::new(), circuit)?;

    let mut file = File::create(dst.join("r1cs"))?;
    r1cs.write_to(&mut file)?;

    let (pk, vk) = Groth16::setup(&r1cs)?;

    let mut file = File::create(dst.join("pk"))?;
    pk.write_raw_to(&mut file)?;

    let mut file = File::create(dst.join("vk"))?;
    vk.write_raw_to(&mut file)?;

    let mut file = File::create(dst.join("Verifier.sol"))?;
    vk.export_solidity(&mut file)?;

    Ok(())
}