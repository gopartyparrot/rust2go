use alloc::vec::Vec;
use std::mem::MaybeUninit;
use std::slice;

/// Returns a string from WebAssembly compatible numeric types representing
/// its pointer and length.
pub unsafe fn ptr_to_string(ptr: u32, len: u32) -> String {
    let slice = slice::from_raw_parts_mut(ptr as *mut u8, len as usize);
    let utf8 = std::str::from_utf8_unchecked_mut(slice);
    return String::from(utf8);
}

/// Returns a pointer and size pair for the given string in a way compatible
/// with WebAssembly numeric types.
///
/// Note: This doesn't change the ownership of the String. To intentionally
/// leak it, use [`std::mem::forget`] on the input after calling this.
pub unsafe fn string_to_ptr(s: &String) -> (u32, u32) {
    return (s.as_ptr() as u32, s.len() as u32);
}

/// Allocates size bytes and leaks the pointer where they start.
pub fn allocate(size: usize) -> *mut u8 {
    // Allocate the amount of bytes needed.
    let vec: Vec<MaybeUninit<u8>> = Vec::with_capacity(size);

    // into_raw leaks the memory to the caller.
    Box::into_raw(vec.into_boxed_slice()) as *mut u8
}

/// Retakes the pointer which allows its memory to be freed.
pub unsafe fn deallocate(ptr: *mut u8, size: usize) {
    let _ = Vec::from_raw_parts(ptr, 0, size);
}

pub unsafe fn process_str<F>(process_fn: F, ptr: u32, len: u32) -> u64
where
    F: Fn(&String) -> String,
{
    let input_str = &ptr_to_string(ptr, len);
    let ret_str = process_fn(input_str);
    let (ptr, len) = string_to_ptr(&ret_str);

    // Note: This changes ownership of the pointer to the external caller. If
    // we didn't call forget, the caller would read back a corrupt value. Since
    // we call forget, the caller must deallocate externally to prevent leaks.
    std::mem::forget(ret_str);
    return ((ptr as u64) << 32) | len as u64;
}
