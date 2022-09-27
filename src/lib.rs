extern crate alloc;
extern crate core;
extern crate wee_alloc;

mod mem;
mod pure;

/// WebAssembly export that accepts a string (linear memory offset, byteCount)
/// and returns a pointer/size pair packed into a u64.
///
/// Note: The return value is leaked to the caller, so it must call
/// [`deallocate`] when finished.
/// Note: This uses a u64 instead of two result values for compatibility with
/// WebAssembly 1.0.
#[cfg_attr(all(target_arch = "wasm32"), export_name = "f64_to_fix_bits")]
#[no_mangle]
pub unsafe extern "C" fn _f64_to_fix_bits(ptr: u32, len: u32) -> u64 {
    // log(&["_f64_to_fix_bits >> ", &mem::ptr_to_string(ptr, len)].concat()); //example debug code
    mem::process_str(pure::f64_to_fix_bits, ptr, len)
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "u128bits_to_fix")]
#[no_mangle]
pub unsafe extern "C" fn _u128bits_to_fix(ptr: u32, len: u32) -> u64 {
    mem::process_str(pure::u128bits_to_fix, ptr, len)
}

/// Set the global allocator to the WebAssembly optimized one.
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

/// WebAssembly export that allocates a pointer (linear memory offset) that can
/// be used for a string.
///
/// This is an ownership transfer, which means the caller must call
/// [`deallocate`] when finished.
#[cfg_attr(all(target_arch = "wasm32"), export_name = "allocate")]
#[no_mangle]
pub extern "C" fn _allocate(size: u32) -> *mut u8 {
    mem::allocate(size as usize)
}

/// WebAssembly export that deallocates a pointer of the given size (linear
/// memory offset, byteCount) allocated by [`allocate`].
#[cfg_attr(all(target_arch = "wasm32"), export_name = "deallocate")]
#[no_mangle]
pub unsafe extern "C" fn _deallocate(ptr: u32, size: u32) {
    mem::deallocate(ptr as *mut u8, size as usize);
}

/// Logs a message to the console using [`_log`]. (don't delete)
fn log(message: &String) {
    unsafe {
        let (ptr, len) = mem::string_to_ptr(message);
        _log(ptr, len);
    }
}

#[link(wasm_import_module = "env")]
extern "C" {
    /// WebAssembly import which prints a string (linear memory offset,
    /// byteCount) to the console.
    ///
    /// Note: This is not an ownership transfer: Rust still owns the pointer
    /// and ensures it isn't deallocated during this call.
    #[link_name = "log"]
    fn _log(ptr: u32, size: u32);
}
