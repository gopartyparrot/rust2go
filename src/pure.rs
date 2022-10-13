use fixed::types::U64F64;

pub fn f64_to_fix_bits(f_str: &String) -> String {
    let f = f_str.parse::<f64>();
    if f.is_err() {
        return format!("ERR: invalid f64 {}", f_str);
    }
    U64F64::from_num(f.unwrap()).to_bits().to_string()
}

pub fn u128bits_to_fix(bits_str: &String) -> String {
    let u = bits_str.parse::<u128>();
    if u.is_err() {
        return format!("ERR: invalid u128 {}", bits_str);
    }
    U64F64::from_bits(u.unwrap()).to_string()
}

#[cfg(test)]
mod tests {

    use crate::pure::*;
    use fixed::types::U64F64;

    #[test]
    fn test_u128_max() {
        println!("{}", u128::MAX);
    }

    #[test]
    fn test_u128_to_fix() {
        println!(
            "{}",
            u128bits_to_fix(&"1341723281558402372940796526592".to_string())
        )
    }
}
