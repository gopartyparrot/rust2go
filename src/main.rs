use clap::Parser;
use std::ffi::CStr;

mod cmd;

fn main() {
    // cli();
    wasi();
}

fn wasi() {
    unsafe {
        let (argc, buf_size) = wasi::args_sizes_get().unwrap();
        let mut argv = Vec::with_capacity(argc);
        let mut buf = Vec::with_capacity(buf_size);
        wasi::args_get(argv.as_mut_ptr(), buf.as_mut_ptr()).unwrap();
        argv.set_len(argc);
        let mut ret = Vec::with_capacity(argc);
        for ptr in argv {
            let s = CStr::from_ptr(ptr.cast());
            // println!("{:?}", s);
            ret.push(s.to_str().unwrap());
        }

        let ret = cmd::Cmd::try_parse_from(ret);
        if ret.is_err() {
            eprint!("{}", ret.unwrap_err().to_string());
            return;
        }
        match ret.unwrap().execute() {
            Ok(data) => print!("{}", data),
            Err(err) => eprint!("{}", err),
        };
    }
}

// fn cli() {
// let cmd = cmd::Cmd::parse();
// print!("{}", cmd.execute());
// }
