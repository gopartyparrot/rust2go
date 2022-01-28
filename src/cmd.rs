use clap::Parser;

// use clap::Parser;
use fixed::types::U64F64;

#[derive(Parser, Debug)]
pub enum Cmd {
    ParseU64F64 {
        #[clap(long, required = true)]
        str: String,
    },
}

impl Cmd {
    pub fn execute(&self) -> Result<String, String> {
        match self {
            Cmd::ParseU64F64 { str } => match str.parse::<u128>() {
                Ok(u) => Ok(U64F64::from_bits(u).to_string()),
                _ => Err(format!("ParseIntError for \"{}\"", str)),
            },
        }
    }
}

#[cfg(test)]
mod tests {
    use fixed::types::U64F64;

    #[test]
    pub fn test_u64f64() {
        let f = U64F64::from_num(123.4);
        println!("{} {}", f.to_bits(), f.to_string())
    }
}
